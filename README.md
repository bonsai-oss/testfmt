# testfmt

Testfmt is a command line tool executing tests in a given directory and formatting the output in common machine-readable formats.

## Usage

All executable files in the given directory will be executed. The output of the tests will be formatted in the given format and written to the output file.
Currently, the following formats are supported:
  - `junit` - JUnit XML format
  - `none` - No formatting, just print the output to stdout

### Gitlab CI

This tool is intended to be used to run integation tests in Gitlab CI. The following example shows how to use it in a `.gitlab-ci.yml` file.
You just need to replace the `./test/integration` path with the path to your integration tests.

```yaml
integration test:
  image: debian:stable-slim
  before_script:
    - wget https://gitlab.com/bonsai-oss/tools/testfmt/-/jobs/5697209953/artifacts/raw/build/testfmt-linux-amd64 -O testfmt && chmod +x testfmt
  script:
    - ./testfmt -f junit -o test-results.xml -d ./test/integration
  artifacts:
    reports:
      junit: test-results.xml
    when: always
```