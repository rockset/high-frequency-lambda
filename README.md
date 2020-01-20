# high-frequency-lambda

If you want to trigger a lambda faster than once per minute, you are out of luck.

This step-function and accompanying iterator lambda lets you invoke your lambda faster than that. This example
triggers it every 10 seconds.

First deploy the lambda you want to trigger, then use this terraform module to deploy the iterator which
will trigger your lambda.

```
module "hfl" {
  source = "github.com/rockset/high-frequency-lambda"
  providers = {
    aws = aws.us-west-2
  }
  target_lambda = "arn of the lambda you want to invoke"
}
```