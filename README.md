# Four Key Metrics

Export four key metrics of your repositories in CSV format.

## Dependencies

* Git

## Installation

* Mac 64-bit: https://github.com/trendyol/metrics/releases/download/v0.1.0/metrics_v0.1.0_darwin_x86_64.tar.gz
* Linux 64-bit: https://github.com/trendyol/metrics/releases/download/v0.1.0/metrics_v0.1.0_linux_arm64.tar.gz
* Windows 64-bit: https://github.com/trendyol/metrics/releases/download/v0.1.0/metrics_v0.1.0_windows_x86_64.zip

## Usage

```bash
Usage: metrics [options...] [repositories...]

Options:
  -r Releases tag regex pattern. Default is: "releases/\\d+"
  -f Fix tag regex pattern. Default is: "releases/fix/\\d+"
  -s Earliest start date of sprint. Date format must be RFC3339. For example: 2020-04-22T16:00:00+03:00
  -d Sprint days size. Default is 7 days

For Example: metrics -s 2020-04-22T16:00:00+03:00 ~/code/my-api ~/code/my-api-2
```

## How To Contribute

Contributions are **welcome** and will be fully **credited**.

Please read the [CONTRIBUTING](CONTRIBUTING.md) and [CODE_OF_CONDUCT](CODE_OF_CONDUCT) files for details.