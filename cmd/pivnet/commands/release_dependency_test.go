package commands_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/go-pivnet"
	"github.com/pivotal-cf-experimental/go-pivnet/cmd/pivnet/commands"
	"github.com/pivotal-cf-experimental/go-pivnet/cmd/pivnet/errorhandler/errorhandlerfakes"
	"github.com/pivotal-cf-experimental/go-pivnet/cmd/pivnet/printer"

	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("release dependency commands", func() {
	var (
		server *ghttp.Server

		fakeErrorHandler *errorhandlerfakes.FakeErrorHandler

		field     reflect.StructField
		outBuffer bytes.Buffer

		productSlug string

		release             pivnet.Release
		releases            []pivnet.Release
		releaseDependencies []pivnet.ReleaseDependency

		responseStatusCode int
		response           interface{}

		releasesResponseStatusCode int
		releasesResponse           pivnet.ReleasesResponse
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		commands.Pivnet.Host = server.URL()

		outBuffer = bytes.Buffer{}
		commands.OutputWriter = &outBuffer
		commands.Printer = printer.NewPrinter(commands.OutputWriter)

		fakeErrorHandler = &errorhandlerfakes.FakeErrorHandler{}
		commands.ErrorHandler = fakeErrorHandler

		productSlug = "some-product-slug"

		release = pivnet.Release{
			ID:      1234,
			Version: "some-release-version",
		}

		releases = []pivnet.Release{
			release,
			{
				ID:      2345,
				Version: "another-release-version",
			},
		}

		releaseDependencies = []pivnet.ReleaseDependency{
			{
				Release: pivnet.DependentRelease{
					ID:      1234,
					Version: "Some version",
				},
			},
			{
				Release: pivnet.DependentRelease{
					ID:      2345,
					Version: "Another version",
				},
			},
		}

		releasesResponseStatusCode = http.StatusOK

		releasesResponse = pivnet.ReleasesResponse{
			Releases: releases,
		}

		responseStatusCode = http.StatusOK
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("ReleasesDependenciesCommand", func() {
		var (
			command commands.ReleaseDependenciesCommand
		)

		BeforeEach(func() {
			responseStatusCode = http.StatusOK

			command = commands.ReleaseDependenciesCommand{
				ProductSlug:    productSlug,
				ReleaseVersion: releases[0].Version,
			}

			response = pivnet.ReleaseDependenciesResponse{
				ReleaseDependencies: releaseDependencies,
			}
		})

		JustBeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fmt.Sprintf("%s/products/%s/releases", apiPrefix, productSlug)),
					ghttp.RespondWithJSONEncoded(releasesResponseStatusCode, releasesResponse),
				),
			)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(
						"GET",
						fmt.Sprintf(
							"%s/products/%s/releases/%d/dependencies",
							apiPrefix,
							productSlug,
							releases[0].ID,
						),
					),
					ghttp.RespondWithJSONEncoded(responseStatusCode, response),
				),
			)
		})

		It("lists all release dependencies for the provided product slug and release version", func() {
			err := command.Execute(nil)
			Expect(err).NotTo(HaveOccurred())

			var returned []pivnet.ReleaseDependency

			err = json.Unmarshal(outBuffer.Bytes(), &returned)
			Expect(err).NotTo(HaveOccurred())

			Expect(returned).To(Equal(releaseDependencies))
		})

		Context("when there is an error", func() {
			BeforeEach(func() {
				responseStatusCode = http.StatusTeapot
			})

			It("invokes the error handler", func() {
				err := command.Execute(nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
			})
		})

		Context("when there is an error getting all releases", func() {
			BeforeEach(func() {
				releasesResponseStatusCode = http.StatusTeapot
			})

			It("invokes the error handler", func() {
				err := command.Execute(nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeErrorHandler.HandleErrorCallCount()).To(Equal(1))
			})
		})

		Describe("ProductSlug flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ReleaseDependenciesCommand{}, "ProductSlug")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("p"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("product-slug"))
			})
		})

		Describe("ReleaseVersion flag", func() {
			BeforeEach(func() {
				field = fieldFor(commands.ReleaseDependenciesCommand{}, "ReleaseVersion")
			})

			It("is required", func() {
				Expect(isRequired(field)).To(BeTrue())
			})

			It("contains short name", func() {
				Expect(shortTag(field)).To(Equal("v"))
			})

			It("contains long name", func() {
				Expect(longTag(field)).To(Equal("release-version"))
			})
		})
	})
})
