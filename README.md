# antrea.io website

Source for the antrea.io website

This repository is a work-in-progess. The source of the current version of the
antrea.io website is located in a
[branch](https://github.com/antrea-io/antrea/tree/website-with-versioned-docs)
of the main Antrea Github repository.

## Building the website locally

### Prerequisites

* [Hugo](https://github.com/gohugoio/hugo)
    * macOS: `brew install hugo`
    * Windows: `choco install hugo-extended -confirm`

### Build

```bash
hugo server --disableFastRender
```

### Access

Access site at http://localhost:1313
