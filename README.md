# LM Prompt Image Tool

This tool takes a prompt and a series of images and returns model outputs as structured, new-lined delimited JSON data which can be used by tools like jq.

## Installation

Make sure you have Go 1.22+ installed then run:

```sh
go install github.com/seanshahkarami/lm-prompt-image@latest
```

You _may_ have to add `$HOME/go/bin/` to your PATH, if you're not seeing the `lm-prompt-image` executable.

## Usage

```sh
lm-prompt-image 'Some prompt' [imagepath1 imagepath2 ...]
```

For example, if I have two images I can view and pipe the results to a `results.ndjson` file.

```sh
$ lm-prompt-image 'Describe this image in detail.' wp3478887.jpg wp9277430.jpg | tee results.ndjson
{"path":"wp3478887.jpg","prompt":"Describe this image in detail.","output":"The scene features a dining area with multiple chairs and tables arranged on the patio space, likely within a restaurant or cafe setting. The patio is illuminated by street lamps at night time, creating an inviting atmosphere for customers. In the image, there are several chairs and dining tables spread throughout the outdoor space. \n\nAdditionally, there are potted plants placed on the premises, adding greenery and enhancing the aesthetic appeal of the area. A handbag can be found on one of the chairs, suggesting that someone might have just arrived at the location or is about to leave. Overall, this scene depicts an inviting outdoor dining space with ample seating arrangements for customers to enjoy their meals in a relaxed and comfortable setting."}
{"path":"wp9277430.jpg","prompt":"Describe this image in detail.","output":"This image features a small bird perched on the branches of a tree with bright green leaves. The bird appears to be looking at something while sitting on top of a twig, possibly observing its surroundings or searching for food. It is captured against a backdrop of a grassy field, making it an interesting and natural scene."}
```

To make this easier to read, I'll use jq to pretty print it:

```sh
jq . results.ndjson 
{
  "path": "wp3478887.jpg",
  "prompt": "Describe this image in detail.",
  "output": "The scene features a dining area with multiple chairs and tables arranged on the patio space, likely within a restaurant or cafe setting. The patio is illuminated by street lamps at night time, creating an inviting atmosphere for customers. In the image, there are several chairs and dining tables spread throughout the outdoor space. \n\nAdditionally, there are potted plants placed on the premises, adding greenery and enhancing the aesthetic appeal of the area. A handbag can be found on one of the chairs, suggesting that someone might have just arrived at the location or is about to leave. Overall, this scene depicts an inviting outdoor dining space with ample seating arrangements for customers to enjoy their meals in a relaxed and comfortable setting."
}
{
  "path": "wp9277430.jpg",
  "prompt": "Describe this image in detail.",
  "output": "This image features a small bird perched on the branches of a tree with bright green leaves. The bird appears to be looking at something while sitting on top of a twig, possibly observing its surroundings or searching for food. It is captured against a backdrop of a grassy field, making it an interesting and natural scene."
}
```

By default, it talks to the model listening on `http://localhost:1234` using the OpenAI chat protocol. This can be configured with the `-addr otherhost:port` flag.

It can also be easily used with a remote llama.cpp server by opening an ssh proxy.

## Advanced Usage

This tool is intended to be using in pipelines. Here's an example of processing an entire directory:

```
$ find path/to/images -name '*.jpg' | xargs lm-prompt-image 'Describe this image in detail.' | tee results.ndjson
```
