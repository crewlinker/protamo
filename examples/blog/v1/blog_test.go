package blogv1_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	blogv1 "github.com/crewlinker/protamo/examples/blog/v1"
	"github.com/crewlinker/protamo/protamoattr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestV1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "examples/blog/v1")
}

// option to ignore exported fields of these types.
var ignoredUnexported = cmpopts.IgnoreUnexported(
	types.AttributeValueMemberS{},
	types.AttributeValueMemberSS{},
	types.AttributeValueMemberNULL{},
	types.AttributeValueMemberM{},
	types.AttributeValueMemberN{},
	types.AttributeValueMemberL{},
	types.AttributeValueMemberBOOL{},
	types.AttributeValueMemberBS{},
)

// option to configure our attribute encoder.
func withJSONTagKey(opts *protamoattr.EncoderOptions) {
	opts.TagKey = "json"
}

var _ = DescribeTable("encoding (with json tag key)", func(msg proto.Message, exp map[string]types.AttributeValue) {
	act, err := protamoattr.MarshalMapWithOptions(msg, withJSONTagKey)
	Expect(err).ToNot(HaveOccurred())

	diff := cmp.Diff(act, exp, ignoredUnexported)

	Expect(diff).To(BeEmpty())
},
	Entry("zero value", &blogv1.BlogPost{}, map[string]types.AttributeValue{}),

	Entry("just scalar field", &blogv1.BlogPost{Title: "foo"}, map[string]types.AttributeValue{
		"1": &types.AttributeValueMemberS{Value: "foo"},
	}),

	Entry("all values", &blogv1.BlogPost{
		Title: "foo",
		Author: &blogv1.BlogAuthor{
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: timestamppb.New(time.Unix(1690280190, 0)),
		},
		Tags: []*blogv1.Tag{
			{Slug: "tutorial", Label: "Tutorial", Color: blogv1.TagColor_TAG_COLOR_BLUE},
		},
		Image: &blogv1.BlogPost_Bitmap{Bitmap: &blogv1.BitmapImage{Src: "/foo.png"}},
		Related: map[string]*blogv1.BlogPost{
			"cool": {Title: "some other blogpost"},
		},
		Metadata: lo.Must1(structpb.NewStruct(map[string]any{
			"foo": 100,
		})),
	}, map[string]types.AttributeValue{
		"1": &types.AttributeValueMemberS{Value: "foo"},
		"2": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"1": &types.AttributeValueMemberS{Value: "John"},
			"2": &types.AttributeValueMemberS{Value: "Doe"},
			"3": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberN{Value: "1690280190"},
			}},
		}},
		"3": &types.AttributeValueMemberL{Value: []types.AttributeValue{
			&types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberS{Value: "tutorial"},
				"2": &types.AttributeValueMemberS{Value: "Tutorial"},
				"3": &types.AttributeValueMemberN{Value: "1"},
			}},
		}},
		"5": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"1": &types.AttributeValueMemberS{Value: "/foo.png"},
		}},
		"6": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"cool": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"1": &types.AttributeValueMemberS{Value: "some other blogpost"},
			}},
		}},
		"7": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"1": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				"foo": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"2": &types.AttributeValueMemberN{Value: "100"},
				}},
			}},
		}},
	}),
)
