import { useState, useEffect } from "react";
import {
  getRecipes,
  getFileURL,
  deleteRecipe,
  getRecipeSteps,
  updateRecipeSteps,
  getFriendlyErrorMessage,
} from "../api";
import { Button } from "../components/Button";
import { Header } from "../components/Header";
import { RecipeSteps } from "../components/RecipeSteps";
import { RecipeStepsEditor } from "../components/RecipeStepsEditor";
import { TagFilter } from "../components/TagFilter";
import type {
  Recipe,
  File,
  RecipeStepResponse as RecipeStep,
} from "../types.gen";
import { asUUID, type Email } from "../branded";
import type { Page } from "../types";

type Props = {
  email: Email;
  currentPage: Page;
  onNavigate: (page: Page) => void;
};

type RecipeWithFiles = Recipe & { files: File[] };

function getParams(): { tags: string[]; recipe: string } {
  const params = new URLSearchParams(window.location.search);
  const tagParam = params.get("tags") ?? "";
  return {
    tags: tagParam ? tagParam.split(",") : [],
    recipe: params.get("recipe") ?? "",
  };
}

function setParams(tags: string[], recipe: string) {
  const params = new URLSearchParams();
  if (tags.length > 0) params.set("tags", tags.join(","));
  if (recipe !== "") params.set("recipe", recipe);
  const search = params.toString();
  const url = search === "" ? "/browse" : `/browse?${search}`;
  window.history.replaceState(null, "", url);
}

export function Browse({ email, currentPage, onNavigate }: Props) {
  const initialParams = getParams();
  const [recipes, setRecipes] = useState<RecipeWithFiles[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedRecipeId, setSelectedRecipeId] = useState<string>(
    initialParams.recipe,
  );
  const [tagFilter, setTagFilter] = useState<string[]>(initialParams.tags);
  const [deleting, setDeleting] = useState(false);
  const [steps, setSteps] = useState<RecipeStep[]>([]);
  const [loadingSteps, setLoadingSteps] = useState(false);
  const [editingSteps, setEditingSteps] = useState(false);
  const [savingSteps, setSavingSteps] = useState(false);

  useEffect(() => {
    getRecipes(email)
      .then((response) => {
        const fileData = response.fileData;
        const recipesWithFiles = response.recipeData.map((recipe) => ({
          ...recipe,
          files: fileData
            .filter((f) => f.recipe_uuid === recipe.uuid)
            .sort((a, b) => a.page_number - b.page_number),
        }));
        setRecipes(recipesWithFiles);
      })
      .catch((err: unknown) => {
        setError(getFriendlyErrorMessage(err, "Failed to load recipes"));
      })
      .finally(() => {
        setLoading(false);
      });
  }, [email]);

  useEffect(() => {
    setParams(tagFilter, selectedRecipeId);
  }, [tagFilter, selectedRecipeId]);

  useEffect(() => {
    if (selectedRecipeId === "") {
      setSteps([]);
      setEditingSteps(false);
      return;
    }
    setLoadingSteps(true);
    setEditingSteps(false);
    getRecipeSteps(asUUID(selectedRecipeId))
      .then((response) => {
        setSteps(response.steps);
      })
      .catch(() => {
        setSteps([]);
      })
      .finally(() => {
        setLoadingSteps(false);
      });
  }, [selectedRecipeId]);

  const handleSaveSteps = async (newSteps: RecipeStep[]) => {
    setSavingSteps(true);
    try {
      await updateRecipeSteps(asUUID(selectedRecipeId), newSteps);
      setSteps(newSteps);
      setEditingSteps(false);
    } catch (err: unknown) {
      setError(getFriendlyErrorMessage(err, "Failed to save steps"));
    } finally {
      setSavingSteps(false);
    }
  };

  const filteredRecipes =
    tagFilter.length === 0
      ? recipes
      : recipes.filter((r) =>
          tagFilter.every((tag) =>
            r.tag_string.toLowerCase().includes(tag.toLowerCase()),
          ),
        );

  const selectedRecipe =
    recipes.find((r) => r.uuid === selectedRecipeId) ?? null;

  const handleSelectRecipe = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedRecipeId(event.target.value);
  };

  const handleDelete = () => {
    if (selectedRecipeId === "") return;
    if (!window.confirm("Are you sure you want to delete this recipe?")) return;

    setDeleting(true);
    deleteRecipe(asUUID(selectedRecipeId))
      .then(() => {
        setRecipes((prev) => prev.filter((r) => r.uuid !== selectedRecipeId));
        setSelectedRecipeId("");
      })
      .catch((err: unknown) => {
        setError(getFriendlyErrorMessage(err, "Failed to delete recipe"));
      })
      .finally(() => {
        setDeleting(false);
      });
  };

  return (
    <>
      <Header email={email} currentPage={currentPage} onNavigate={onNavigate} />
      <div className="container">
        <h1>Browse Recipes</h1>

        {error !== null && <p className="error">{error}</p>}

        <div className="flex-row" style={{ marginBottom: "1rem" }}>
          <TagFilter value={tagFilter} onChange={setTagFilter} />
        </div>

        {loading ? (
          <p>Loading...</p>
        ) : (
          <select
            className="select"
            onChange={handleSelectRecipe}
            value={selectedRecipeId}
          >
            <option value="" disabled>
              Select a recipe
            </option>
            {filteredRecipes.map((recipe) => (
              <option key={recipe.uuid} value={recipe.uuid}>
                {recipe.name}
                {recipe.tag_string ? ` - ${recipe.tag_string}` : ""}
              </option>
            ))}
          </select>
        )}

        {selectedRecipe !== null && (
          <div style={{ marginTop: "2rem" }}>
            <div
              className="flex-row"
              style={{ alignItems: "center", gap: "1rem" }}
            >
              <h2 style={{ margin: 0 }}>{selectedRecipe.name}</h2>
              <Button
                variant="danger"
                onClick={handleDelete}
                disabled={deleting}
              >
                {deleting ? "Deleting..." : "Delete"}
              </Button>
            </div>
            <p>Tags: {selectedRecipe.tag_string}</p>
            <div
              className="flex-row"
              style={{ alignItems: "center", gap: "1rem", marginTop: "1.5rem" }}
            >
              <h3 style={{ margin: 0 }}>Steps</h3>
              {!editingSteps && !loadingSteps && (
                <Button
                  onClick={() => {
                    setEditingSteps(true);
                  }}
                >
                  Edit
                </Button>
              )}
            </div>
            {loadingSteps ? (
              <p>Loading steps...</p>
            ) : editingSteps ? (
              <RecipeStepsEditor
                steps={steps}
                onSave={handleSaveSteps}
                onCancel={() => {
                  setEditingSteps(false);
                }}
                saving={savingSteps}
              />
            ) : (
              <RecipeSteps steps={steps} />
            )}
            {selectedRecipe.files.map((file) => (
              <img
                key={file.uuid}
                src={getFileURL(file.url)}
                alt={`${selectedRecipe.name} page ${String(file.page_number + 1)}`}
                className="recipe-image"
              />
            ))}
          </div>
        )}
      </div>
    </>
  );
}
