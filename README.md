# rigel

Small &amp; fast image proxy server written in Go!

- Blazingly fast.
- Wide range of image operations (e.g., Resize, Enlarge, Crop [including Smart Crop], Rotate, Flip, Flop, Zoom, Thumbnail, Extract area, Watermark [text or image], Gaussian blur effect, Custom output color space [RGB, grayscale...], Format conversion [with additional quality/compression settings], EXIF metadata [size, alpha channel, profile, orientation...], Trim).
- Precaching with heads-up requests, providing short URLs.
- Multi-level caching (memory(lfu), redis, etc) .
- Signed URLs with expiry support.

## Getting started

### Installation

It's already dockerized, Just run:

```bash
docker run -p 8080:8080 sdqr/rigel
```

Otherwise, you can fork source code and customize it for your needs.

### How to use

- Check [rigelsdk](https://github.com/sdqri/rigelsdk).

## Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](issues).
