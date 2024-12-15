"use client";

import type { PutBlobResult } from "@vercel/blob";
// import { setMaxIdleHTTPParsers } from "http";
import { useState, useRef } from "react";

export default function AvatarUploadPage() {
  const inputFileRef = useRef<HTMLInputElement>(null);
  const metadataRef = useRef<HTMLInputElement>(null);
  const filenameRef = useRef<HTMLInputElement>(null);
  const [imageBlob, setImageBlob] = useState<PutBlobResult | null>(null);
  const [metadataBlob, setMetadataBlob] = useState<PutBlobResult | null>(null);
  console.log(metadataBlob, setMetadataBlob);
  return (
    <>
      <h1>Upload Your Avatar</h1>

      <form
        onSubmit={async (event) => {
          event.preventDefault();
          await sendUpload(
            inputFileRef,
            setImageBlob,
            metadataRef,
            setMetadataBlob,
            filenameRef
          );
        }}
      >
        <input
          name="file"
          ref={inputFileRef}
          type="file"
          required
          className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <input
          name="metadata"
          ref={metadataRef}
          type="text"
          placeholder="Enter a list of tags, separated by commas"
          className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
        />
        <input
          name="filename"
          ref={filenameRef}
          type="text"
          placeholder="Enter your desired file name"
          className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
        />
        <button
          type="submit"
          className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          Upload
        </button>
      </form>
      {imageBlob && (
        <div>
          Blob url: <a href={imageBlob.url}>{imageBlob.url}</a>
        </div>
      )}
    </>
  );
}
async function sendUpload(
  inputFileRef: React.RefObject<HTMLInputElement | null>,
  setImageBlob: React.Dispatch<React.SetStateAction<PutBlobResult | null>>,
  metadataRef: React.RefObject<HTMLInputElement | null>,
  setMetedataBlob: React.Dispatch<React.SetStateAction<PutBlobResult | null>>,
  filenameRef: React.RefObject<HTMLInputElement | null>
) {
  if (!inputFileRef.current?.files) {
    throw new Error("No file selected");
  }
  const file = inputFileRef.current.files[0];
  if (!metadataRef.current?.value) {
    metadataRef.current = null;
  }
  const tagString = metadataRef.current?.value;
  let filename = file.name;
  if (filenameRef.current?.value) {
    filename = filenameRef.current.value;
  }

  const url = `/api/recipe/upload?filename=${filename}&tagString=${tagString}`;

  console.log("inpage, sending url: ", url);
  const response = await fetch(url, {
    method: "POST",
    body: file,
  });
  console.log("got response, ", response);

  const newBlob = (await response.json()) as PutBlobResult;
  console.log("newBlob", JSON.stringify(newBlob));
  console.log(setMetedataBlob);
  setImageBlob(newBlob);
}

/*
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
  const response = await fetch(`/api/recipe/upload?filename=${file.name}`, {
    method: "POST",
    body: file,
  });

  const newBlob = (await response.json()) as PutBlobResult;

  setImageBlob(newBlob);* /
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

  setMetedataBlob(newBlob);* /
}
*/
