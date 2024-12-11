"use client";

import type { PutBlobResult } from "@vercel/blob";
// import { setMaxIdleHTTPParsers } from "http";
import { useState, useRef } from "react";

export default function AvatarUploadPage() {
  const inputFileRef = useRef<HTMLInputElement>(null);
  const metadataRef = useRef<HTMLInputElement>(null);
  const [imageBlob, setImageBlob] = useState<PutBlobResult | null>(null);
  const [metadataBlob, setMetadataBlob] = useState<PutBlobResult | null>(null);
  console.log(metadataBlob, setMetadataBlob);
  return (
    <>
      <h1>Upload Your Avatar</h1>

      <form
        onSubmit={async (event) => {
          event.preventDefault();
          await sendImage(inputFileRef, setImageBlob);
          await sendMetadata(metadataRef, setMetadataBlob, imageBlob);
        }}
      >
        <input name="file" ref={inputFileRef} type="file" required />
        <input
          name="metadata"
          ref={metadataRef}
          type="text"
          placeholder="Enter a list of tags, separated by commas"
        />
        <button type="submit">Upload</button>
      </form>
      {imageBlob && (
        <div>
          Blob url: <a href={imageBlob.url}>{imageBlob.url}</a>
        </div>
      )}
    </>
  );
}

async function sendImage(
  inputFileRef: React.RefObject<HTMLInputElement | null>,
  setImageBlob: React.Dispatch<React.SetStateAction<PutBlobResult | null>>
) {
  if (!inputFileRef.current?.files) {
    throw new Error("No file selected");
  }
  console.log("setImageBlob", setImageBlob);
  const file = inputFileRef.current.files[0];
  console.log("file", file);
  /*
  const response = await fetch(`/api/avatar/upload?filename=${file.name}`, {
    method: "POST",
    body: file,
  });

  const newBlob = (await response.json()) as PutBlobResult;

  setImageBlob(newBlob);*/
}

async function sendMetadata(
  metadataRef: React.RefObject<HTMLInputElement | null>,
  setMetedataBlob: React.Dispatch<React.SetStateAction<PutBlobResult | null>>,
  filenameBlob: PutBlobResult | null
) {
  if (!metadataRef.current) {
    console.log("No metadata current");
    return;
  }
  const tags = metadataRef.current.value.split(", ");
  console.log("tags:", tags);
  if (!metadataRef.current?.files) {
    console.log("No metadata file selected");
    return;
  }
  return;
  /*  const file = metadataRef.current.files[0];

  const response = await fetch(`/api/avatar/upload?filename=${filename}`, {
    method: "POST",
    body: file,
  });

  const newBlob = (await response.json()) as PutBlobResult;

  setMetedataBlob(newBlob);*/
}
