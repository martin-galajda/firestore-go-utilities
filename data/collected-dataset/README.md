# Documentation for the collected dataset

The whole collected dataset for annotating important images in the websites has been exported into JSON.

The structure of the JSON is described below:
```
[
  {
    Metadata: {
      AnnotatedElementsData: {
        "<website_url>": {
          "DataAnnotationId": " --- unique ID generated for annotating the image and inserteds into HTML document as data-annotation-id attribute",
          "ElemPathFromRoot": " --s- path from the root element of the HTML document",
          "ImgUrlBase64": " --- base64 encoding for the URL of the image",
          "Url": " --- URL of the image"
        }
      },
      Html: " --- whole HTML document, with 'data-annotation-id' marks"
    }
    Url: "--- URL of the website",
  }

]

```
