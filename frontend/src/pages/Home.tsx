type Props = {
  onNavigate: (page: 'upload' | 'browse') => void
}

export function Home({ onNavigate }: Props) {
  return (
    <div className="center">
      <h1>Welcome to Hungr!</h1>
      <p>Would you like to upload a recipe or browse recipes?</p>
      <div className="flex-row">
        <button className="btn" onClick={() => { onNavigate('upload') }}>
          Upload Recipe
        </button>
        <button className="btn" onClick={() => { onNavigate('browse') }}>
          Browse Recipes
        </button>
      </div>
    </div>
  )
}
