package pivnet_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logger"
	"github.com/pivotal-cf/go-pivnet/logger/loggerfakes"
)

var _ = Describe("PivnetClient - release upgrade paths", func() {
	var (
		server     *ghttp.Server
		client     pivnet.Client
		token      string
		apiAddress string
		userAgent  string

		newClientConfig pivnet.ClientConfig
		fakeLogger      logger.Logger

		releaseID int
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		apiAddress = server.URL()
		token = "my-auth-token"
		userAgent = "pivnet-resource/0.1.0 (some-url)"

		releaseID = 2345

		fakeLogger = &loggerfakes.FakeLogger{}
		newClientConfig = pivnet.ClientConfig{
			Host:      apiAddress,
			Token:     token,
			UserAgent: userAgent,
		}
		client = pivnet.NewClient(newClientConfig, fakeLogger)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Get", func() {
		It("returns the release upgrade paths", func() {
			response := pivnet.ReleaseUpgradePathsResponse{
				ReleaseUpgradePaths: []pivnet.ReleaseUpgradePath{
					{
						Release: pivnet.UpgradePathRelease{
							ID:      9876,
							Version: "release 9876",
						},
					},
					{
						Release: pivnet.UpgradePathRelease{
							ID:      8765,
							Version: "release 8765",
						},
					},
				},
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf(
						"%s/products/%s/releases/%d/upgrade_paths",
						apiPrefix,
						productSlug,
						releaseID,
					)),
					ghttp.RespondWithJSONEncoded(http.StatusOK, response),
				),
			)

			releaseUpgradePaths, err := client.ReleaseUpgradePaths.Get(productSlug, releaseID)
			Expect(err).NotTo(HaveOccurred())

			Expect(releaseUpgradePaths).To(HaveLen(2))
			Expect(releaseUpgradePaths[0].Release.ID).To(Equal(9876))
			Expect(releaseUpgradePaths[1].Release.ID).To(Equal(8765))
		})

		Context("when the server responds with a non-2XX status code", func() {
			var (
				body []byte
			)

			BeforeEach(func() {
				body = []byte(`{"message":"foo message"}`)
			})

			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", fmt.Sprintf(
							"%s/products/%s/releases/%d/upgrade_paths",
							apiPrefix,
							productSlug,
							releaseID,
						)),
						ghttp.RespondWith(http.StatusTeapot, body),
					),
				)
			})

			It("returns an error", func() {
				_, err := client.ReleaseUpgradePaths.Get(productSlug, releaseID)
				Expect(err.Error()).To(ContainSubstring("foo message"))
			})
		})
	})

	Describe("Add", func() {
		var (
			previousReleaseID int
		)

		BeforeEach(func() {
			previousReleaseID = 1234
		})

		It("adds the release upgrade path", func() {
			expectedRequestBody := `{"upgrade_path":{"release_id":1234}}`

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", fmt.Sprintf(
						"%s/products/%s/releases/%d/add_upgrade_path",
						apiPrefix,
						productSlug,
						releaseID,
					)),
					ghttp.VerifyJSON(expectedRequestBody),
					ghttp.RespondWithJSONEncoded(http.StatusNoContent, nil),
				),
			)

			err := client.ReleaseUpgradePaths.Add(
				productSlug,
				releaseID,
				previousReleaseID,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the server responds with a non-2XX status code", func() {
			var (
				body []byte
			)

			BeforeEach(func() {
				body = []byte(`{"message":"foo message"}`)
			})

			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PATCH", fmt.Sprintf(
							"%s/products/%s/releases/%d/add_upgrade_path",
							apiPrefix,
							productSlug,
							releaseID,
						)),
						ghttp.RespondWith(http.StatusTeapot, body),
					),
				)
			})

			It("returns an error", func() {
				err := client.ReleaseUpgradePaths.Add(
					productSlug,
					releaseID,
					previousReleaseID,
				)
				Expect(err.Error()).To(ContainSubstring("foo message"))
			})
		})
	})

	Describe("Remove", func() {
		var (
			previousReleaseID int
		)

		BeforeEach(func() {
			previousReleaseID = 1234
		})

		It("removes the release upgrade path", func() {
			expectedRequestBody := `{"upgrade_path":{"release_id":1234}}`

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", fmt.Sprintf(
						"%s/products/%s/releases/%d/remove_upgrade_path",
						apiPrefix,
						productSlug,
						releaseID,
					)),
					ghttp.VerifyJSON(expectedRequestBody),
					ghttp.RespondWithJSONEncoded(http.StatusNoContent, nil),
				),
			)

			err := client.ReleaseUpgradePaths.Remove(
				productSlug,
				releaseID,
				previousReleaseID,
			)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the server responds with a non-2XX status code", func() {
			var (
				body []byte
			)

			BeforeEach(func() {
				body = []byte(`{"message":"foo message"}`)
			})

			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PATCH", fmt.Sprintf(
							"%s/products/%s/releases/%d/remove_upgrade_path",
							apiPrefix,
							productSlug,
							releaseID,
						)),
						ghttp.RespondWith(http.StatusTeapot, body),
					),
				)
			})

			It("returns an error", func() {
				err := client.ReleaseUpgradePaths.Remove(
					productSlug,
					releaseID,
					previousReleaseID,
				)
				Expect(err.Error()).To(ContainSubstring("foo message"))
			})
		})
	})
})
