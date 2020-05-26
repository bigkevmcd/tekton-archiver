# tekton-archiver

Watches PipelineRuns and TaskRuns and archives the output upon completion.

## Installation

A [Deployment](./deploy) is provided, you can modify to build your own image.

Out of the box, this uses the S3 archiver to archive logs.


```shell
$ kubectl create secret generic archiver-secret --from-literal=aws-access-key-id=$AWS_SECRET_ACCESS_KEY --from-literal=aws-secret-access-key=$AWS_SECRET_ACCESS_KEY
```

```
$ kubectl create configmap archiver-config --from-literal=bucket-name=$AWS_BUCKET --from-literal=aws-region=$AWS_REGION
```
