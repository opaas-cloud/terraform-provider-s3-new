package macie2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/macie2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfmacie2 "github.com/hashicorp/terraform-provider-aws/internal/service/macie2"
)

func testAccClassificationJob_basic(t *testing.T) {
	var macie2Output macie2.DescribeClassificationJobOutput
	resourceName := "aws_macie2_classification_job.test"
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClassificationJobDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, macie2.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccClassificationJobConfig_nameGenerated(bucketName, macie2.JobTypeOneTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output),
					acctest.CheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "terraform-"),
					resource.TestCheckResourceAttr(resourceName, "job_type", macie2.JobTypeOneTime),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.bucket_definitions.0.buckets.0", bucketName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClassificationJob_Name_Generated(t *testing.T) {
	var macie2Output macie2.DescribeClassificationJobOutput
	resourceName := "aws_macie2_classification_job.test"
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClassificationJobDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, macie2.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccClassificationJobConfig_nameGenerated(bucketName, macie2.JobTypeOneTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output),
					acctest.CheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "terraform-"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClassificationJob_NamePrefix(t *testing.T) {
	var macie2Output macie2.DescribeClassificationJobOutput
	resourceName := "aws_macie2_classification_job.test"
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	namePrefix := "tf-acc-test-prefix-"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClassificationJobDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, macie2.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccClassificationJobConfig_namePrefix(bucketName, namePrefix, macie2.JobTypeOneTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "name", namePrefix),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", namePrefix),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClassificationJob_disappears(t *testing.T) {
	var macie2Output macie2.DescribeClassificationJobOutput
	resourceName := "aws_macie2_classification_job.test"
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClassificationJobDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, macie2.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccClassificationJobConfig_nameGenerated(bucketName, macie2.JobTypeOneTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output),
					acctest.CheckResourceDisappears(acctest.Provider, tfmacie2.ResourceClassificationJob(), resourceName),
				),
			},
		},
	})
}

func testAccClassificationJob_Status(t *testing.T) {
	var macie2Output, macie2Output2 macie2.DescribeClassificationJobOutput
	resourceName := "aws_macie2_classification_job.test"
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClassificationJobDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, macie2.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccClassificationJobConfig_status(bucketName, macie2.JobStatusRunning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output),
					resource.TestCheckResourceAttr(resourceName, "job_type", macie2.JobTypeScheduled),
					resource.TestCheckResourceAttr(resourceName, "job_status", macie2.JobStatusRunning),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.bucket_definitions.0.buckets.0", bucketName),
					acctest.CheckResourceAttrRFC3339(resourceName, "created_at"),
				),
			},
			{
				Config: testAccClassificationJobConfig_status(bucketName, macie2.JobStatusUserPaused),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output2),
					testAccCheckClassificationJobNotRecreated(&macie2Output, &macie2Output2),
					resource.TestCheckResourceAttr(resourceName, "job_type", macie2.JobTypeScheduled),
					resource.TestCheckResourceAttr(resourceName, "job_status", macie2.JobStatusUserPaused),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.bucket_definitions.0.buckets.0", bucketName),
					acctest.CheckResourceAttrRFC3339(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClassificationJob_complete(t *testing.T) {
	var macie2Output macie2.DescribeClassificationJobOutput
	resourceName := "aws_macie2_classification_job.test"
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	description := "Description of a test"
	descriptionUpdated := "Updated Description of a test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClassificationJobDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, macie2.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccClassificationJobConfig_complete(bucketName, macie2.JobStatusRunning, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output),
					resource.TestCheckResourceAttr(resourceName, "job_type", macie2.JobTypeScheduled),
					resource.TestCheckResourceAttr(resourceName, "job_status", macie2.JobStatusRunning),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.bucket_definitions.0.buckets.0", bucketName),
					acctest.CheckResourceAttrRFC3339(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "sampling_percentage", "100"),
					resource.TestCheckResourceAttr(resourceName, "initial_run", "true"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.excludes.0.and.0.simple_scope_term.0.comparator", "EQ"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.excludes.0.and.0.simple_scope_term.0.key", "OBJECT_EXTENSION"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.excludes.0.and.0.simple_scope_term.0.values.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.includes.0.and.0.simple_scope_term.0.comparator", "EQ"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.includes.0.and.0.simple_scope_term.0.key", "OBJECT_EXTENSION"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.includes.0.and.0.simple_scope_term.0.values.0", "test"),
				),
			},
			{
				Config: testAccClassificationJobConfig_complete(bucketName, macie2.JobStatusRunning, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output),
					resource.TestCheckResourceAttr(resourceName, "job_type", macie2.JobTypeScheduled),
					resource.TestCheckResourceAttr(resourceName, "job_status", macie2.JobStatusRunning),
					acctest.CheckResourceAttrRFC3339(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
					resource.TestCheckResourceAttr(resourceName, "sampling_percentage", "100"),
					resource.TestCheckResourceAttr(resourceName, "initial_run", "true"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.excludes.0.and.0.simple_scope_term.0.comparator", "EQ"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.excludes.0.and.0.simple_scope_term.0.key", "OBJECT_EXTENSION"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.excludes.0.and.0.simple_scope_term.0.values.0", "test"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.includes.0.and.0.simple_scope_term.0.comparator", "EQ"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.includes.0.and.0.simple_scope_term.0.key", "OBJECT_EXTENSION"),
					resource.TestCheckResourceAttr(resourceName, "s3_job_definition.0.scoping.0.includes.0.and.0.simple_scope_term.0.values.0", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClassificationJob_WithTags(t *testing.T) {
	var macie2Output macie2.DescribeClassificationJobOutput
	resourceName := "aws_macie2_classification_job.test"
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClassificationJobDestroy,
		ErrorCheck:               acctest.ErrorCheck(t, macie2.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccClassificationJobConfig_completeTags(bucketName, macie2.JobStatusRunning),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output),
					resource.TestCheckResourceAttr(resourceName, "job_type", macie2.JobTypeScheduled),
					resource.TestCheckResourceAttr(resourceName, "job_status", macie2.JobStatusRunning),
					acctest.CheckResourceAttrRFC3339(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Key", "value"),
				),
			},
			{
				Config: testAccClassificationJobConfig_completeTags(bucketName, macie2.JobStatusUserPaused),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClassificationJobExists(resourceName, &macie2Output),
					resource.TestCheckResourceAttr(resourceName, "job_type", macie2.JobTypeScheduled),
					resource.TestCheckResourceAttr(resourceName, "job_status", macie2.JobStatusUserPaused),
					acctest.CheckResourceAttrRFC3339(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Key", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClassificationJobExists(resourceName string, macie2Session *macie2.DescribeClassificationJobOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).Macie2Conn
		input := &macie2.DescribeClassificationJobInput{JobId: aws.String(rs.Primary.ID)}

		resp, err := conn.DescribeClassificationJob(input)

		if err != nil {
			return err
		}

		if resp == nil {
			return fmt.Errorf("macie ClassificationJob %q does not exist", rs.Primary.ID)
		}

		*macie2Session = *resp

		return nil
	}
}

func testAccCheckClassificationJobDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).Macie2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_macie2_classification_job" {
			continue
		}
		input := &macie2.DescribeClassificationJobInput{JobId: aws.String(rs.Primary.ID)}

		resp, err := conn.DescribeClassificationJob(input)

		if tfawserr.ErrMessageContains(err, macie2.ErrCodeAccessDeniedException, "Macie is not enabled") ||
			tfawserr.ErrMessageContains(err, macie2.ErrCodeValidationException, "cannot update cancelled job for job") {
			continue
		}

		if err != nil {
			return err
		}

		if resp != nil && aws.StringValue(resp.JobStatus) != macie2.JobStatusCancelled {
			return fmt.Errorf("macie ClassificationJob %q still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckClassificationJobNotRecreated(i, j *macie2.DescribeClassificationJobOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !aws.TimeValue(i.CreatedAt).Equal(aws.TimeValue(j.CreatedAt)) {
			return fmt.Errorf("Macie Classification Job recreated")
		}

		return nil
	}
}

func testAccClassificationJobConfig_nameGenerated(bucketName, jobType string) string {
	return fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_macie2_account" "test" {}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_macie2_classification_job" "test" {
  depends_on = [aws_macie2_account.test]
  job_type   = %[2]q
  s3_job_definition {
    bucket_definitions {
      account_id = data.aws_caller_identity.current.account_id
      buckets    = [aws_s3_bucket.test.bucket]
    }
  }
}
`, bucketName, jobType)
}

func testAccClassificationJobConfig_namePrefix(nameBucket, namePrefix, jobType string) string {
	return fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_macie2_account" "test" {}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_macie2_classification_job" "test" {
  name_prefix = %[2]q
  job_type    = %[3]q
  s3_job_definition {
    bucket_definitions {
      account_id = data.aws_caller_identity.current.account_id
      buckets    = [aws_s3_bucket.test.bucket]
    }
  }
  depends_on = [aws_macie2_account.test]
}
`, nameBucket, namePrefix, jobType)
}

func testAccClassificationJobConfig_complete(nameBucket, jobStatus, description string) string {
	return fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_macie2_account" "test" {}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_macie2_classification_job" "test" {
  job_type = "SCHEDULED"
  s3_job_definition {
    bucket_definitions {
      account_id = data.aws_caller_identity.current.account_id
      buckets    = [aws_s3_bucket.test.bucket]
    }
    scoping {
      excludes {
        and {
          simple_scope_term {
            comparator = "EQ"
            key        = "OBJECT_EXTENSION"
            values     = ["test"]
          }
        }
      }
      includes {
        and {
          simple_scope_term {
            comparator = "EQ"
            key        = "OBJECT_EXTENSION"
            values     = ["test"]
          }
        }
      }
    }
  }
  schedule_frequency {
    daily_schedule = true
  }
  sampling_percentage = 100
  initial_run         = true
  job_status          = %[2]q
  description         = %[3]q

  depends_on = [aws_macie2_account.test]
}
`, nameBucket, jobStatus, description)
}

func testAccClassificationJobConfig_status(nameBucket, jobStatus string) string {
	return fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_macie2_account" "test" {}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_macie2_classification_job" "test" {
  job_type = "SCHEDULED"
  s3_job_definition {
    bucket_definitions {
      account_id = data.aws_caller_identity.current.account_id
      buckets    = [aws_s3_bucket.test.bucket]
    }
    scoping {
      excludes {
        and {
          simple_scope_term {
            comparator = "EQ"
            key        = "OBJECT_EXTENSION"
            values     = ["test"]
          }
        }
      }
      includes {
        and {
          simple_scope_term {
            comparator = "EQ"
            key        = "OBJECT_EXTENSION"
            values     = ["test"]
          }
        }
      }
    }
  }
  schedule_frequency {
    daily_schedule = true
  }
  sampling_percentage = 100
  initial_run         = true
  job_status          = %[2]q

  depends_on = [aws_macie2_account.test]
}
`, nameBucket, jobStatus)
}

func testAccClassificationJobConfig_completeTags(nameBucket, jobStatus string) string {
	return fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_macie2_account" "test" {}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_macie2_classification_job" "test" {
  job_type = "SCHEDULED"
  s3_job_definition {
    bucket_definitions {
      account_id = data.aws_caller_identity.current.account_id
      buckets    = [aws_s3_bucket.test.bucket]
    }
    scoping {
      excludes {
        and {
          simple_scope_term {
            comparator = "EQ"
            key        = "OBJECT_EXTENSION"
            values     = ["test"]
          }
        }
      }
      includes {
        and {
          simple_scope_term {
            comparator = "EQ"
            key        = "OBJECT_EXTENSION"
            values     = ["test"]
          }
        }
      }
    }
  }
  schedule_frequency {
    daily_schedule = true
  }
  sampling_percentage = 100
  description         = "test"
  initial_run         = true
  job_status          = %[2]q
  tags = {
    Key = "value"
  }

  depends_on = [aws_macie2_account.test]
}
`, nameBucket, jobStatus)
}
