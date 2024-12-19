"use client";
import { useState, useEffect } from "react";
import Image from "next/image";

const USERID = 1; // Temporarily keeping this simple

type Metadata = {
  filename: string;
  tagString: string;
  createdAt: string;
  imageUrl: string;
};

function packageData(data: object[]): Metadata[] {
  return data.map((item) => {
    const validItem = item as {
      filename: string;
      tag_string: string;
      created_at: string;
      url: string;
    };

    if (
      typeof validItem !== "object" ||
      !validItem ||
      typeof validItem.filename !== "string" ||
      typeof validItem.tag_string !== "string" ||
      typeof validItem.created_at !== "string" ||
      typeof validItem.url !== "string"
    ) {
      throw new Error("Data is missing required fields");
    }

    return {
      filename: validItem.filename,
      tagString: validItem.tag_string,
      createdAt: validItem.created_at,
      imageUrl: validItem.url,
    };
  });
}

async function fetchMetadata(userId: number): Promise<Metadata[]> {
  const url = `/api/recipe/upload?type=metadata&user_id=${userId}`;
  const response = await fetch(url);
  if (!response.ok) throw new Error("Failed to fetch data");
  const result = await response.json();
  return packageData(result);
}

async function fetchImageDetails(imageUrl: string): Promise<{
  image: Blob;
  dimensions: { width: number; height: number };
}> {
  const response = await fetch(imageUrl);
  if (!response.ok) throw new Error("Failed to fetch image");

  const result = await response.blob();
  const img = new window.Image();
  const url = URL.createObjectURL(result);
  img.src = url;

  return new Promise((resolve) => {
    img.onload = () => {
      resolve({
        image: result,
        dimensions: { width: img.naturalWidth, height: img.naturalHeight },
      });
    };
  });
}

export default function ShowRecipe() {
  const [data, setData] = useState<Metadata[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [selectedOption, setSelectedOption] = useState<string>("");
  const [image, setImage] = useState<Blob | null>(null);
  const [imageDims, setImageDims] = useState<{
    width: number;
    height: number;
  } | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  const userId = USERID;

  console.log("Entering showrecipe");

  useEffect(() => {
    (async () => {
      try {
        setLoading(true);
        const metadata = await fetchMetadata(userId);
        setData(metadata);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Unknown error occurred");
      } finally {
        setLoading(false);
      }
    })();
  }, [userId]);

  const handleSelectionChange = async (imageUrl: string) => {
    try {
      const { image, dimensions } = await fetchImageDetails(imageUrl);
      setImage(image);
      setImageDims(dimensions);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error occurred");
    }
  };

  return (
    <div>
      <h1>Please choose the recipe you want to load.</h1>
      {error && <div>Error: {error}</div>}
      {loading ? (
        <div>
          <h2>Loading...</h2>
          <select
            className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled
          >
            <option>Loading options...</option>
          </select>
        </div>
      ) : (
        data &&
        data.length > 0 && (
          <div>
            <h2>Data:</h2>
            <select
              value={selectedOption}
              onChange={(event) => {
                const selectedValue = event.target.value;
                setSelectedOption(selectedValue);
                handleSelectionChange(selectedValue);
              }}
              className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="" disabled>
                Select an option
              </option>
              {data.map((item, index) => (
                <option key={index} value={item.imageUrl}>
                  {item.filename} - {item.tagString}
                </option>
              ))}
            </select>
          </div>
        )
      )}
      {image && imageDims && (
        <div>
          <h2>Image:</h2>
          <Image
            src={URL.createObjectURL(image)}
            alt="Fetched from database"
            width={imageDims.width}
            height={imageDims.height}
          />
        </div>
      )}
    </div>
  );
}
