# go-tex2svg

> Golang HTTP LaTeX compiler server

## How it works

- [fasthttp](https://github.com/fasthttp) serve a HTTP server with a POST `text/plain` route at `/latex` 
- saves a `.tex` file
- [pdflatex](https://linux.die.net/man/1/pdflatex) compiles `.tex` and output `.pdf`
- [pdf2svg](https://github.com/dawbarton/pdf2svg) converts `.pdf` to `.svg`
- [minify](https://github.com/tdewolff/minify) minifies `.svg`
- finally, responds SVG content with `image/svg+xml` header

## Usage

### Build 

```bash
docker build -t=go-tex2svg .
```

### Run

```bash
docker run -d -p 4000:4000 tex2svg
```

### Convert

```
curl -d "$\\frac{1}{2}$" http://localhost:4000/latex
```

<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="6.09pt" height="16.543pt" viewBox="0 0 6.09 16.543"><defs><g><symbol overflow="visible" id="glyph0-0"><path style="stroke:none;" d="M2.390625-1.265625v-.125c0-.8125 1.203125-1.390625 1.203125-2.5C3.59375-4.65625 2.9375-5 2.203125-5 1.453125-5 .9375-4.578125.9375-3.96875c0 .359375.140625.53125.390625.53125.171875.0.328125-.109375.328125-.328125C1.65625-4.0625 1.5-4.0625 1.5-4.3125c0-.234375.203125-.40625.640625-.40625.484375.0.78125.265625.78125.859375.0 1.140625-.828125 1.421875-.828125 2.390625v.203125zm.28125.90625c0-.234375-.15625-.4375-.421875-.4375s-.4375.203125-.4375.4375c0 .25.140625.421875.4375.421875.265625.0.421875-.1875.421875-.421875zm1.765625 1.625v-7.125H.078125v7.125zm-.375-.40625H.453125v-6.34375H4.0625zm0 0"/></symbol><symbol overflow="visible" id="glyph0-1"><path style="stroke:none;" d="M3.1875.0V-.203125c-.6875.0-.96875-.140625-.96875-.5v-4.0625H2l-1.453125.40625v.25c.234375-.078125.625-.125.765625-.125.1875.0.25.109375.25.359375V-.703125c0 .359375-.265625.5-.96875.5V0zm0 0"/></symbol><symbol overflow="visible" id="glyph0-2"><path style="stroke:none;" d="M3.46875-1.1875H3.25c-.171875.4375-.265625.59375-.671875.59375H1.546875l-.53125.03125V-.609375l1-.921875c.78125-.8125 1.203125-1.328125 1.203125-1.984375.0-.78125-.578125-1.296875-1.421875-1.296875-.734375.0-1.234375.421875-1.46875 1.140625l.1875.078125c.265625-.5625.59375-.734375 1.078125-.734375.578125.0.9375.359375.9375.921875.0.78125-.375 1.25-1.109375 2.015625L.28125-.25V0h2.9375zm0 0"/></symbol></g><clipPath id="clip1"><path d="M1 11H5v5.542969H1zm0 0"/></clipPath></defs><g id="surface1"><g style="fill:rgb(0%,0%,0%);fill-opacity:1;"><use xlink:href="#glyph0-1" x="1.195" y="4.887"/></g><path style="fill:none;stroke-width:0.677;stroke-linecap:butt;stroke-linejoin:miter;stroke:rgb(0%,0%,0%);stroke-opacity:1;stroke-miterlimit:10;" d="M3125e-7-53125e-8H3.699531" transform="matrix(1,0,0,-1,1.195,8.144)"/><g clip-path="url(#clip1)" clip-rule="nonzero"><g style="fill:rgb(0%,0%,0%);fill-opacity:1;"><use xlink:href="#glyph0-2" x="1.195" y="16.543"/></g></g></g></svg>

```
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="6.09pt" height="16.543pt" viewBox="0 0 6.09 16.543"><defs><g><symbol overflow="visible" id="glyph0-0"><path style="stroke:none;" d="M2.390625-1.265625v-.125c0-.8125 1.203125-1.390625 1.203125-2.5C3.59375-4.65625 2.9375-5 2.203125-5 1.453125-5 .9375-4.578125.9375-3.96875c0 .359375.140625.53125.390625.53125.171875.0.328125-.109375.328125-.328125C1.65625-4.0625 1.5-4.0625 1.5-4.3125c0-.234375.203125-.40625.640625-.40625.484375.0.78125.265625.78125.859375.0 1.140625-.828125 1.421875-.828125 2.390625v.203125zm.28125.90625c0-.234375-.15625-.4375-.421875-.4375s-.4375.203125-.4375.4375c0 .25.140625.421875.4375.421875.265625.0.421875-.1875.421875-.421875zm1.765625 1.625v-7.125H.078125v7.125zm-.375-.40625H.453125v-6.34375H4.0625zm0 0"/></symbol><symbol overflow="visible" id="glyph0-1"><path style="stroke:none;" d="M3.1875.0V-.203125c-.6875.0-.96875-.140625-.96875-.5v-4.0625H2l-1.453125.40625v.25c.234375-.078125.625-.125.765625-.125.1875.0.25.109375.25.359375V-.703125c0 .359375-.265625.5-.96875.5V0zm0 0"/></symbol><symbol overflow="visible" id="glyph0-2"><path style="stroke:none;" d="M3.46875-1.1875H3.25c-.171875.4375-.265625.59375-.671875.59375H1.546875l-.53125.03125V-.609375l1-.921875c.78125-.8125 1.203125-1.328125 1.203125-1.984375.0-.78125-.578125-1.296875-1.421875-1.296875-.734375.0-1.234375.421875-1.46875 1.140625l.1875.078125c.265625-.5625.59375-.734375 1.078125-.734375.578125.0.9375.359375.9375.921875.0.78125-.375 1.25-1.109375 2.015625L.28125-.25V0h2.9375zm0 0"/></symbol></g><clipPath id="clip1"><path d="M1 11H5v5.542969H1zm0 0"/></clipPath></defs><g id="surface1"><g style="fill:rgb(0%,0%,0%);fill-opacity:1;"><use xlink:href="#glyph0-1" x="1.195" y="4.887"/></g><path style="fill:none;stroke-width:0.677;stroke-linecap:butt;stroke-linejoin:miter;stroke:rgb(0%,0%,0%);stroke-opacity:1;stroke-miterlimit:10;" d="M3125e-7-53125e-8H3.699531" transform="matrix(1,0,0,-1,1.195,8.144)"/><g clip-path="url(#clip1)" clip-rule="nonzero"><g style="fill:rgb(0%,0%,0%);fill-opacity:1;"><use xlink:href="#glyph0-2" x="1.195" y="16.543"/></g></g></g></svg>
```