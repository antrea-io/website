# antrea.io website

This repository holds the source code for generating the
[antrea.io](https://antrea.io) website. The documentation contents for the
website are not primarily hosted here, but instead are mirrored from the main
Antrea [repository](https://github.com/antrea-io/antrea). To improve or fix a
documentation page, you will need to open a PR in that repository, and not this
one.

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

## Removing old versions

When a version of Antrea is no longer supported, it is a good idea to remove it
from the drop-down navigation menu, by editing `config.yaml`. The actual
documentation does not need to be removed from the `contents/` directory, unless
website size becomes an issue. In doing so, we ensure that old links will keep
working, while also decluttering the website.
