# tekton-archiver

**This is pre-alpha code**

Watches PipelineRuns and archives the output upon completion, currently only to AWS S3.

## Installation

A [Deployment](./deploy) is provided, you can modify to build your own image.

Out of the box, this uses the S3 Archiver to archive logs, to use this, you'll
need to create a secret with your credentials.

If you want to keep your bucket and region out of the secret, you'll need to edit the default deployment.yaml.

```shell
$ kubectl create secret generic archiver-secret --from-literal=AWS_ACCESS_KEY_ID=$AWS_SECRET_ACCESS_KEY --from-literal=AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY --from-literal=BUCKET_NAME=$AWS_BUCKET --from-literal=AWS_REGION=$AWS_REGION
```

## Labelling PipelineRuns

Currently, the archiver only archives logs that have a label with
`"pipelinerun": "archive"`

```yaml
apiVersion: tekton.dev/v1alpha1
kind: PipelineRun
metadata:
  name: demo-pipelinerun
  labels:
    "pipelinerun": "archive"
spec:
```


## Testing

```shell
$ go test ./...
```
