# testfmt

Testfmt is a command line tool executing tests in a given directory and formatting the output in common machine-readable formats.

## Usage

### Gitlab CI

```yaml
integration test:
  script:
    - testfmt -f junit -o test-results.xml -d ./test/integration
  artifacts:
    reports:
      junit: test-results.xml
    when: always
```