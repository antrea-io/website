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

## Website updates

The website is automatically updated every time the `antrea-io/antrea` `main`
branch is updated (or rather, every time the documentation is updated), and
every time a new Antrea release is created. When either one of these events
happens, the `Update website source` Github
[workflow](.github/workflows/update.yml) is triggered in this repo. The
necessary scripts take care of updating the correct docs, with the changes being
committed to `main`. by the Github workflow directly (as the `antrea-bot` Github
user). The workflow takes care of updating all the necessary "metadata" in the
case of a new Antrea release: TOC file, TOC mapping file, Hugo's configuration
file, etc.

To modify website source files, you can open a PR by following the Antrea
contribution
[guidelines](https://github.com/antrea-io/antrea/blob/main/CONTRIBUTING.md). Please
note that manual changes to [content/docs/](content/docs/) should be avoided. In
particular, any manual change to the `main` version of the docs will be
overwriten by the `Update website source` Github workflow.

## Define a new Golang vanity import path

To configure the vanity import path for a new Golang module hosted in the
antrea-io Github organization, you need to open a PR with the following changes:

 * add the appropriate entry to `static/_redirects`
 * add a new HTML file under `static/golang/`, e.g., `static/golang/<new repo
   name>.html`; add the necessary HTML content (look at existing files for
   reference)
