package grpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"

	"io/ioutil"

	"time"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-function-discovery/internal/grpc"
	"github.com/solo-io/gloo-plugins/grpc"
	"github.com/solo-io/gloo-storage/dependencies/file"
)

var _ = Describe("Discovergrpc", func() {
	Describe("happy path", func() {
		Context("upstream for a grpc server", func() {
			It("returns service info for grpc", func() {
				dir, err := ioutil.TempDir("", "")
				Expect(err).To(BeNil())
				files, err := file.NewFileStorage(dir, time.Millisecond)
				Expect(err).To(BeNil())
				detector := NewGRPCDetector(files)
				addr := fmt.Sprintf("localhost:%v", port)
				svcInfo, annotations, err := detector.DetectFunctionalService(&v1.Upstream{Name: "Test"}, addr)
				Expect(err).To(BeNil())
				Expect(annotations).To(BeNil())
				fileRef := "grpc-discovery:Bookstore.descriptors"
				Expect(svcInfo).To(Equal(&v1.ServiceInfo{
					Type: grpc.ServiceTypeGRPC,
					Properties: grpc.EncodeServiceProperties(grpc.ServiceProperties{
						GRPCServiceNames:   []string{"Bookstore"},
						DescriptorsFileRef: fileRef,
					}),
				}))
				list, err := files.List()
				Expect(err).To(BeNil())
				Expect(list).To(HaveLen(1))
				Expect(list[0].Ref).To(Equal(fileRef))
			})
		})
	})
})
