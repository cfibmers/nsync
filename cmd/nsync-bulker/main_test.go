package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	"github.com/cloudfoundry-incubator/bbs"
	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/cloudfoundry-incubator/locket"
	"github.com/cloudfoundry-incubator/runtime-schema/cc_messages"
)

var _ = Describe("Syncing desired state with CC", func() {
	const interruptTimeout = 5 * time.Second

	var (
		fakeCC *ghttp.Server

		process ifrit.Process

		domainTTL time.Duration

		bulkerLockName    = "nsync_bulker_lock"
		pollingInterval   time.Duration
		heartbeatInterval time.Duration

		logger lager.Logger
	)

	startBulker := func(check bool) ifrit.Process {
		runner := ginkgomon.New(ginkgomon.Config{
			Name:          "nsync-bulker",
			AnsiColorCode: "97m",
			StartCheck:    "nsync.bulker.started",
			Command: exec.Command(
				bulkerPath,
				"-ccBaseURL", fakeCC.URL(),
				"-pollingInterval", pollingInterval.String(),
				"-domainTTL", domainTTL.String(),
				"-bulkBatchSize", "10",
				"-lifecycle", "buildpack/some-stack:some-health-check.tar.gz",
				"-lifecycle", "docker:the/docker/lifecycle/path.tgz",
				"-fileServerURL", "http://file-server.com",
				"-lockRetryInterval", "1s",
				"-consulCluster", consulRunner.ConsulCluster(),
				"-bbsAddress", fakeBBS.URL(),
				"-privilegedContainers", "false",
			),
		})

		if !check {
			runner.StartCheck = ""
		}

		return ginkgomon.Invoke(runner)
	}

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")

		fakeCC = ghttp.NewServer()

		pollingInterval = 500 * time.Millisecond
		domainTTL = 1 * time.Second
		heartbeatInterval = 30 * time.Second

		desiredAppResponses := map[string]string{
			"process-guid-1": `{
					"disk_mb": 1024,
					"environment": [
						{ "name": "env-key-1", "value": "env-value-1" },
						{ "name": "env-key-2", "value": "env-value-2" }
					],
					"file_descriptors": 16,
					"num_instances": 42,
					"log_guid": "log-guid-1",
					"memory_mb": 256,
					"process_guid": "process-guid-1",
					"routing_info": {
						"http_routes":
							[
								{"hostname": "route-1"},
								{"hostname": "route-2"},
								{"hostname": "new-route"}
							]
					},
					"droplet_uri": "source-url-1",
					"stack": "some-stack",
					"start_command": "start-command-1",
					"execution_metadata": "execution-metadata-1",
					"health_check_timeout_in_seconds": 123456,
					"etag": "1.1"
				}`,
			"process-guid-2": `{
					"disk_mb": 1024,
					"environment": [
						{ "name": "env-key-1", "value": "env-value-1" },
						{ "name": "env-key-2", "value": "env-value-2" }
					],
					"file_descriptors": 16,
					"num_instances": 4,
					"log_guid": "log-guid-1",
					"memory_mb": 256,
					"process_guid": "process-guid-2",
					"routing_info": {
						"http_routes": [
								{ "hostname": "route-3", "route_service_url":"https://rs.example.com"}
							],
						"tcp_routes": [
							  {"router_group_guid": "guid-1", "external_port":5222, "container_port":60000}
					 	  ]
						},
					"droplet_uri": "source-url-1",
					"stack": "some-stack",
					"start_command": "start-command-1",
					"execution_metadata": "execution-metadata-1",
					"health_check_timeout_in_seconds": 123456,
					"etag": "2.1"
				}`,
			"process-guid-3": `{
					"disk_mb": 512,
					"environment": [],
					"file_descriptors": 8,
					"num_instances": 4,
					"log_guid": "log-guid-3",
					"memory_mb": 128,
					"process_guid": "process-guid-3",
					"routing_info": { "http_routes": [] },
					"droplet_uri": "source-url-3",
					"stack": "some-stack",
					"start_command": "start-command-3",
					"execution_metadata": "execution-metadata-3",
					"health_check_timeout_in_seconds": 123456,
					"etag": "3.1"
				}`,
		}

		fakeCC.RouteToHandler("GET", "/internal/bulk/apps",
			ghttp.RespondWith(200, `{
					"token": {},
					"fingerprints": [
							{
								"process_guid": "process-guid-1",
								"etag": "1.1"
							},
							{
								"process_guid": "process-guid-2",
								"etag": "2.1"
							},
							{
								"process_guid": "process-guid-3",
								"etag": "3.1"
							}
					]
				}`),
		)

		fakeCC.RouteToHandler("POST", "/internal/bulk/apps",
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				var processGuids []string
				decoder := json.NewDecoder(req.Body)
				err := decoder.Decode(&processGuids)
				Expect(err).NotTo(HaveOccurred())

				appResponses := make([]json.RawMessage, 0, len(processGuids))
				for _, processGuid := range processGuids {
					appResponses = append(appResponses, json.RawMessage(desiredAppResponses[processGuid]))
				}

				payload, err := json.Marshal(appResponses)
				Expect(err).NotTo(HaveOccurred())

				w.Write(payload)
			}),
		)
	})

	AfterEach(func() {
		defer fakeCC.Close()
	})

	Describe("when the CC polling interval elapses", func() {
		JustBeforeEach(func() {
			process = startBulker(true)
		})

		AfterEach(func() {
			ginkgomon.Interrupt(process, interruptTimeout)
		})

		Context("once the state has been synced with CC", func() {
			Context("lrps", func() {
				BeforeEach(func() {
					schedulingInfoResponse := models.DesiredLRPSchedulingInfosResponse{
						Error: nil,
						DesiredLrpSchedulingInfos: []*models.DesiredLRPSchedulingInfo{
							{
								DesiredLRPKey: models.DesiredLRPKey{ // perfect. love it. keep it.
									ProcessGuid: "process-guid-1",
									Domain:      cc_messages.AppLRPDomain,
								},
								Annotation: "1.1",
							},
							{
								DesiredLRPKey: models.DesiredLRPKey{ // annotation mismatch so update
									ProcessGuid: "process-guid-2",
									Domain:      cc_messages.AppLRPDomain,
								},
							}, // missing 3 so create it
							{
								DesiredLRPKey: models.DesiredLRPKey{ // extra to be removed
									ProcessGuid: "process-guid-4",
									Domain:      cc_messages.AppLRPDomain,
								},
								Annotation: "4.1",
							},
						},
					}
					data, err := schedulingInfoResponse.Marshal()
					Expect(err).ToNot(HaveOccurred())

					fakeBBS.RouteToHandler("POST", "/v1/desired_lrp_scheduling_infos/list",
						ghttp.RespondWith(200, data, http.Header{bbs.ContentTypeHeader: []string{bbs.ProtoContentType}}),
					)

					fakeBBS.RouteToHandler("POST", "/v1/domains/upsert",
						ghttp.RespondWith(200, `{}`),
					)

					fakeBBS.RouteToHandler("POST", "/v1/desired_lrp/desire",
						ghttp.CombineHandlers(
							ghttp.VerifyContentType("application/x-protobuf"),
							func(w http.ResponseWriter, req *http.Request) {
								body, err := ioutil.ReadAll(req.Body)
								Expect(err).ShouldNot(HaveOccurred())
								defer req.Body.Close()

								protoMessage := &models.DesireLRPRequest{}

								err = proto.Unmarshal(body, protoMessage)
								Expect(err).ToNot(HaveOccurred(), "Failed to unmarshal protobuf")

								Expect(protoMessage.DesiredLrp.ProcessGuid).To(Equal("process-guid-3"))
							},
						),
					)

					fakeBBS.RouteToHandler("POST", "/v1/desired_lrp/update",
						ghttp.CombineHandlers(
							ghttp.VerifyContentType("application/x-protobuf"),
							func(w http.ResponseWriter, req *http.Request) {
								body, err := ioutil.ReadAll(req.Body)
								Expect(err).ShouldNot(HaveOccurred())
								defer req.Body.Close()

								protoMessage := &models.UpdateDesiredLRPRequest{}

								err = proto.Unmarshal(body, protoMessage)
								Expect(err).ToNot(HaveOccurred(), "Failed to unmarshal protobuf")

								Expect(*protoMessage.Update.Annotation).To(Equal("2.1"))
								Expect(protoMessage.ProcessGuid).To(Equal("process-guid-2"))
							},
						),
					)

					expectedLRPDeleteRequest := &models.RemoveDesiredLRPRequest{ProcessGuid: "process-guid-4"}
					fakeBBS.RouteToHandler("POST", "/v1/desired_lrp/remove",
						ghttp.VerifyProtoRepresenting(expectedLRPDeleteRequest),
					)

					fakeCC.RouteToHandler("GET", "/internal/v3/bulk/task_states",
						ghttp.RespondWith(200, `{"token": {},"task_states": []}`),
					)

					fakeBBS.RouteToHandler("POST", "/v1/tasks/list.r1",
						ghttp.RespondWith(200, `{"error": {},"tasks": []}`),
					)
				})

				It("it (adds), (updates), and (removes extra) LRPs", func() {
					Eventually(func() bool {
						for _, r := range fakeBBS.ReceivedRequests() {
							if r.URL.Path == "/v1/desired_lrp/desire" {
								return true
							}
						}
						return false
					}).Should(BeTrue())

					Eventually(func() bool {
						for _, r := range fakeBBS.ReceivedRequests() {
							if r.URL.Path == "/v1/desired_lrp/update" {
								return true
							}
						}
						return false
					}).Should(BeTrue())

					Eventually(func() bool {
						for _, r := range fakeBBS.ReceivedRequests() {
							if r.URL.Path == "/v1/desired_lrp/remove" {
								return true
							}
						}
						return false
					}).Should(BeTrue())
				})
			})

			Context("tasks", func() {
				Context("CC has a task, but the bbs does not", func() {
					BeforeEach(func() {
						fakeCC.RouteToHandler("GET", "/internal/v3/bulk/task_states",
							ghttp.RespondWith(200, `{
							"token": {},
							"task_states": [
								{
									"task_guid": "task-guid-1",
									"state": "RUNNING",
									"completion_callback": "`+fmt.Sprintf("%s/internal/v3/tasks/task-guid-1/completed", fakeCC.URL())+`"
								}
							]
						}`),
						)

						fakeCC.RouteToHandler("POST", "/internal/v3/tasks/task-guid-1/completed",
							ghttp.CombineHandlers(
								ghttp.VerifyJSON(`{
								"task_guid": "task-guid-1",
								"failed": true,
								"failure_reason": "Unable to determine completion status"
							}`),
								ghttp.RespondWith(200, `{}`),
							),
						)

						fakeBBS.RouteToHandler("POST", "/v1/tasks/list.r1",
							ghttp.RespondWith(200, `{"error": {},"tasks": []}`),
						)

						fakeBBS.RouteToHandler("POST", "/v1/desired_lrp_scheduling_infos/list",
							ghttp.RespondWith(200, `{"error":{},"desired_lrp_scheduling_infos":	[]}`),
						)

						fakeBBS.RouteToHandler("POST", "/v1/domains/upsert",
							ghttp.RespondWith(200, `{}`),
						)

						fakeBBS.RouteToHandler("POST", "/v1/desired_lrp/desire",
							ghttp.RespondWith(200, `{}`),
						)
					})

					It("completes the tasks and sets the state to failed", func() {
						var request *http.Request
						Eventually(func() *http.Request {
							for _, r := range fakeCC.ReceivedRequests() {
								if r.URL.Path == "/internal/v3/tasks/task-guid-1/completed" {
									request = r
									return r
								}
							}
							return nil
						}, 2*domainTTL).ShouldNot(BeNil())
					})
				})

				Context("The BBS has a task, but the CC does not", func() {
					BeforeEach(func() {
						fakeCC.RouteToHandler("GET", "/internal/v3/bulk/task_states",
							ghttp.RespondWith(200, `{ "token": {}, "task_states": [] }`),
						)

						taskResponse := models.TasksResponse{
							Tasks: []*models.Task{
								{
									TaskGuid: "task-guid-1",
									State:    models.Task_Completed,
									Domain:   cc_messages.RunningTaskDomain,
									TaskDefinition: &models.TaskDefinition{
										CompletionCallbackUrl: "/internal/v3/tasks/task-guid-1/completed",
									},
								},
							},
						}
						data, err := taskResponse.Marshal()
						Expect(err).ToNot(HaveOccurred())

						fakeBBS.RouteToHandler("POST", "/v1/tasks/list.r1",
							ghttp.RespondWith(200, data, http.Header{bbs.ContentTypeHeader: []string{bbs.ProtoContentType}}),
						)

						fakeBBS.RouteToHandler("POST", "/v1/desired_lrp_scheduling_infos/list",
							ghttp.RespondWith(200, `{"error":{},"desired_lrp_scheduling_infos":	[]}`),
						)

						fakeBBS.RouteToHandler("POST", "/v1/domains/upsert",
							ghttp.RespondWith(200, `{}`),
						)

						fakeBBS.RouteToHandler("POST", "/v1/desired_lrp/desire",
							ghttp.RespondWith(200, `{}`),
						)
					})

					It("cancels the tasks in the BBS", func() {
						Eventually(func() bool {
							for _, r := range fakeBBS.ReceivedRequests() {
								if r.URL.Path == "/v1/tasks/cancel" {
									return true
								}
							}
							return false
						}, 20*domainTTL).ShouldNot(BeNil())
					})
				})
			})

			Describe("domains", func() {
				var (
					foundTaskDomain chan bool
				)

				BeforeEach(func() {
					foundTaskDomain = make(chan bool, 2)

					fakeBBS.RouteToHandler("POST", "/v1/desired_lrp_scheduling_infos/list",
						ghttp.RespondWith(200, `{"error":{},"desired_lrp_scheduling_infos":	[]}`),
					)

					taskResponse := models.TasksResponse{
						Tasks: []*models.Task{
							{
								TaskGuid: "task-guid-1",
								State:    models.Task_Completed,
								Domain:   cc_messages.RunningTaskDomain,
								TaskDefinition: &models.TaskDefinition{
									CompletionCallbackUrl: "/internal/v3/tasks/task-guid-1/completed",
								},
							},
						},
					}
					data, err := taskResponse.Marshal()
					Expect(err).ToNot(HaveOccurred())

					fakeBBS.RouteToHandler("POST", "/v1/tasks/list.r1",
						ghttp.RespondWith(200, data, http.Header{bbs.ContentTypeHeader: []string{bbs.ProtoContentType}}),
					)

					fakeCC.RouteToHandler("GET", "/internal/v3/bulk/task_states",
						ghttp.RespondWith(200, `{
							"token": {},
							"task_states": [
								{
									"task_guid": "task-guid-1",
									"state": "RUNNING",
									"completion_callback": "`+fmt.Sprintf("%s/internal/v3/tasks/task-guid-1/completed", fakeCC.URL())+`"
								}
							]
						}`),
					)

					fakeCC.RouteToHandler("POST", "/internal/v3/tasks/task-guid-1/completed",
						ghttp.CombineHandlers(
							ghttp.VerifyJSON(`{
								"task_guid": "task-guid-1",
								"failed": true,
								"failure_reason": "Unable to determine completion status"
							}`),
							ghttp.RespondWith(200, `{}`),
						),
					)

					fakeBBS.RouteToHandler("POST", "/v1/desired_lrp/desire",
						ghttp.RespondWith(200, `{}`),
					)

					fakeBBS.RouteToHandler("POST", "/v1/tasks/list.r1",
						ghttp.RespondWith(200, `{"error": {},"tasks": []}`),
					)

					fakeBBS.RouteToHandler("POST", "/v1/domains/upsert",
						ghttp.CombineHandlers(
							ghttp.VerifyContentType("application/x-protobuf"),
							func(w http.ResponseWriter, req *http.Request) {
								body, err := ioutil.ReadAll(req.Body)
								Expect(err).ShouldNot(HaveOccurred())
								defer req.Body.Close()

								protoMessage := &models.UpsertDomainRequest{}

								err = proto.Unmarshal(body, protoMessage)
								Expect(err).ToNot(HaveOccurred(), "Failed to unmarshal protobuf")

								if protoMessage.Domain == cc_messages.RunningTaskDomain {
									close(foundTaskDomain)
								}
							},
						),
					)
				})

				It("updates the domains", func() {
					Eventually(foundTaskDomain, 2*domainTTL).Should(BeClosed())
				})
			})
		})
	})

	Context("when the bulker loses the lock", func() {
		BeforeEach(func() {
			heartbeatInterval = 1 * time.Second

			fakeCC.RouteToHandler("GET", "/internal/v3/bulk/task_states",
				ghttp.RespondWith(200, `{"token": {},"task_states": []}`),
			)

			fakeBBS.RouteToHandler("POST", "/v1/tasks/list.r1",
				ghttp.RespondWith(200, `{"error": {},"tasks": []}`),
			)

			fakeBBS.RouteToHandler("POST", "/v1/desired_lrp_scheduling_infos/list",
				ghttp.RespondWith(200, `{"error":{},"desired_lrp_scheduling_infos":	[]}`),
			)

			fakeBBS.RouteToHandler("POST", "/v1/domains/upsert",
				ghttp.RespondWith(200, `{}`),
			)

			fakeBBS.RouteToHandler("POST", "/v1/desired_lrp/desire",
				ghttp.RespondWith(200, `{}`),
			)
		})

		JustBeforeEach(func() {
			process = startBulker(true)

			_, err := consulRunner.NewClient().KV().DeleteTree(locket.LockSchemaPath(bulkerLockName), nil)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			ginkgomon.Interrupt(process, interruptTimeout)
		})

		It("exits with an error", func() {
			Eventually(process.Wait(), 5*domainTTL).Should(Receive(HaveOccurred()))
		})
	})

	Context("when the bulker initially does not have the lock", func() {
		var nsyncLockClaimerProcess ifrit.Process

		BeforeEach(func() {
			heartbeatInterval = 1 * time.Second

			nsyncLockClaimer := locket.NewLock(logger, consulRunner.NewClient(), locket.LockSchemaPath(bulkerLockName), []byte("something-else"), clock.NewClock(), locket.RetryInterval, locket.LockTTL)
			nsyncLockClaimerProcess = ifrit.Invoke(nsyncLockClaimer)
		})

		JustBeforeEach(func() {
			process = startBulker(false)
		})

		AfterEach(func() {
			ginkgomon.Interrupt(process, interruptTimeout)
			ginkgomon.Kill(nsyncLockClaimerProcess)
		})

		It("does not make any requests", func() {
			Consistently(func() int {
				return len(fakeBBS.ReceivedRequests())
			}).Should(Equal(0))
		})

		Context("when the lock becomes available", func() {
			var (
				expectedLRPDomainRequest  *models.UpsertDomainRequest
				expectedTaskDomainRequest *models.UpsertDomainRequest

				foundLRPDomain  chan bool
				foundTaskDomain chan bool
			)

			BeforeEach(func() {
				foundLRPDomain = make(chan bool, 2)
				foundTaskDomain = make(chan bool, 2)

				expectedLRPDomainRequest = &models.UpsertDomainRequest{
					Domain: "cf-app",
					Ttl:    uint32(domainTTL.Seconds()),
				}

				expectedTaskDomainRequest = &models.UpsertDomainRequest{
					Domain: "cf-tasks",
					Ttl:    uint32(domainTTL.Seconds()),
				}

				fakeCC.RouteToHandler("GET", "/internal/v3/bulk/task_states",
					ghttp.RespondWith(200, `{"token": {},"task_states": []}`),
				)

				fakeBBS.RouteToHandler("POST", "/v1/tasks/list.r1",
					ghttp.RespondWith(200, `{"error": {},"tasks": []}`),
				)

				fakeBBS.RouteToHandler("POST", "/v1/desired_lrp_scheduling_infos/list",
					ghttp.RespondWith(200, `{"error":{},"desired_lrp_scheduling_infos":	[]}`),
				)

				fakeBBS.RouteToHandler("POST", "/v1/desired_lrp/desire",
					ghttp.RespondWith(200, `{}`),
				)

				fakeBBS.RouteToHandler("POST", "/v1/domains/upsert",
					ghttp.CombineHandlers(
						ghttp.VerifyContentType("application/x-protobuf"),
						func(w http.ResponseWriter, req *http.Request) {
							body, err := ioutil.ReadAll(req.Body)
							Expect(err).ShouldNot(HaveOccurred())
							defer req.Body.Close()

							protoMessage := &models.UpsertDomainRequest{}

							err = proto.Unmarshal(body, protoMessage)
							Expect(err).ToNot(HaveOccurred(), "Failed to unmarshal protobuf")

							if protoMessage.Domain == cc_messages.AppLRPDomain {
								close(foundLRPDomain)
							}

							if protoMessage.Domain == cc_messages.RunningTaskDomain {
								close(foundTaskDomain)
							}
						},
					),
				)

				ginkgomon.Kill(nsyncLockClaimerProcess)
				time.Sleep(pollingInterval + 10*time.Millisecond)
			})

			It("is updated", func() {
				Eventually(foundLRPDomain, 2*domainTTL).Should(BeClosed())
				Eventually(foundTaskDomain, 2*domainTTL).Should(BeClosed())
			})
		})
	})
})
