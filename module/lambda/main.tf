/*
resource "null_resource" "one" {
  provisioner "local-exec" {
    working_dir = var.filename
    command     = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -o bootstrap"
  }
  triggers = {
    always_run = "${timestamp()}"
  }
}
*/

resource "terraform_data" "go_build" {
  triggers_replace = [timestamp()]

  provisioner "local-exec" {
    environment = {
      GOOS        = "linux"
      GOARCH      = "amd64"
      CGO_ENABLED = "0"

    }
    working_dir = var.filename
    command     = "go build -tags lambda.norpc -o bootstrap"
  }
}

data "archive_file" "archive" {
  depends_on  = [terraform_data.go_build]
  output_path = "${var.filename}.zip"
  source_dir  = var.filename
  type        = "zip"
}

resource "aws_lambda_function" "one" {
  filename         = data.archive_file.archive.output_path
  function_name    = var.function_name
  role             = var.role
  handler          = "hanlder"
  runtime          = "provided.al2"
  source_code_hash = data.archive_file.archive.output_base64sha256
  environment {
    variables = var.environment
  }
}
