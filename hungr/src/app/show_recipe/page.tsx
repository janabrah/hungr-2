"use client";
import { useState, useEffect } from "react";
import Image from "next/image";

const USERID = 1; // Temporarily keeping this simple

type Metadata = {
  filename: string;
  tagString: string;
  createdAt: string;
  imageUrls: string[];
};

function packageData(
  recipeData: object[],
  fileData: object[],
  mappingData: object[]
): Metadata[] {
  console.log(recipeData);
  console.log(fileData);
  console.log(mappingData);
  return recipeData.map((item) => {
    const validItem = item as {
      id: number;
      filename: string;
      tag_string: string;
      created_at: string;
    };

    if (
      typeof validItem !== "object" ||
      !validItem ||
      typeof validItem.filename !== "string" ||
      typeof validItem.tag_string !== "string" ||
      typeof validItem.created_at !== "string"
    ) {
      throw new Error("Data is missing required fields");
    }
    console.log("getting fileIds");
    const fileIds = mappingData.reduce((acc: number[], mappingItem) => {
      const validMappingItem = mappingItem as {
        recipe_id: number;
        file_id: number;
      };
      if (
        typeof validMappingItem !== "object" ||
        !validMappingItem ||
        typeof validMappingItem.recipe_id !== "number"
      ) {
        throw new Error("Data is missing required fields");
      }
      if (validMappingItem.recipe_id === validItem.id) {
        acc.push(validMappingItem.file_id);
      }
      return acc;
    }, [] as number[]);
    console.log("fileIds are: " + fileIds);
    console.log("getting fileUrls");
    const fileUrls = fileData.reduce((acc: string[], fileItem) => {
      const validFileItem = fileItem as { id: number; url: string };
      if (
        typeof validFileItem !== "object" ||
        !validFileItem ||
        typeof validFileItem.url !== "string"
      ) {
        throw new Error("Data is missing required fields");
      }
      if (fileIds.includes(validFileItem.id)) {
        acc.push(validFileItem.url);
      }
      return acc;
    }, [] as string[]);
    console.log("fileUrls is: " + JSON.stringify(fileUrls));
    return {
      filename: validItem.filename,
      tagString: validItem.tag_string,
      createdAt: validItem.created_at,
      imageUrls: fileUrls,
    };
  });
}

async function fetchMetadata(userId: number): Promise<Metadata[]> {
  const url = `/api/recipe/upload?type=metadata&user_id=${userId}`;
  const response = await fetch(url);
  if (!response.ok) throw new Error("Failed to fetch data");
  const result = await response.json();
  console.log("got result: " + JSON.stringify(result));
  return packageData(result.recipeData, result.fileData, result.mappingData);
}

function getTagMetadata(data: Metadata[] | null, tag: string): Metadata[] {
  // When the metadata dataset gets bigger, should probably do a new sql query
  // instead of reprocessing the tag strings
  if (!data) {
    return [];
  }
  console.log("getting tag metadata");
  return data.filter((item) => item.tagString.split(", ").includes(tag));
}

async function fetchImageDetails(imageUrl: string): Promise<{
  image: Blob;
  dimensions: { width: number; height: number };
}> {
  console.log("fetching image: " + imageUrl);
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
  const [selectedOption, setSelectedOption] = useState<string[]>([]);
  const [images, setImages] = useState<
    Array<{ image: Blob; dimensions: { width: number; height: number } }>
  >([]);
  const [loading, setLoading] = useState<boolean>(true);

  const userId = USERID;

  console.log("Entering showrecipe");

  useEffect(() => {
    (async () => {
      try {
        setLoading(true);
        console.log("fetching metadata");
        const metadata = await fetchMetadata(userId);
        console.log("metadata is: " + JSON.stringify(metadata));
        setData(metadata);
      } catch (err) {
        console.log(JSON.stringify(err));
        setError(err instanceof Error ? err.message : "Unknown error occurred");
      } finally {
        setLoading(false);
      }
    })();
  }, [userId]);

  const handleSelectionChange = async (imageUrls: string[]) => {
    try {
      const images = [];
      for (const imageUrl of imageUrls) {
        const { image, dimensions } = await fetchImageDetails(imageUrl);
        images.push({ image, dimensions });
      }
      setImages(images);
    } catch (err) {
      console.log(JSON.stringify(err));
      setError(err instanceof Error ? err.message : "Unknown error occurred");
    }
  };

  const tagSearch = async (event: React.FormEvent) => {
    event.preventDefault();
    const tag = event.currentTarget.querySelector("input")?.value;
    if (!tag) {
      return;
    }
    try {
      setLoading(true);
      setData(getTagMetadata(data, tag));
    } catch (err) {
      console.log(JSON.stringify(err));
      setError(err instanceof Error ? err.message : "Unknown error occurred");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1>If you would like to search for a specific tag, enter it here</h1>
      <form onSubmit={tagSearch}>
        <div className="flex space-x-2">
          <input
            type="text"
            placeholder="Enter a search tag"
            className="block p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
          />
          <button
            type="submit"
            className="block p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Search
          </button>
        </div>
      </form>
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
                console.log("selected value: " + selectedValue);
                setSelectedOption(JSON.parse(selectedValue));
                handleSelectionChange(JSON.parse(selectedValue));
              }}
              className="block w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
            >
              <option value="" disabled>
                Select an option
              </option>
              {data.map((item, index) => (
                <option key={index} value={JSON.stringify(item.imageUrls)}>
                  {item.filename} - {item.tagString}
                </option>
              ))}
            </select>
          </div>
        )
      )}
      {images.length > 0 && (
        <div>
          <h2>Images:</h2>
          {images.map((imgData, index) => (
            <Image
              key={index}
              src={URL.createObjectURL(imgData.image)}
              alt={`Fetched from database ${index}`}
              width={imgData.dimensions.width}
              height={imgData.dimensions.height}
            />
          ))}
        </div>
      )}
    </div>
  );
}
