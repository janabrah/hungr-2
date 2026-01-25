type LoadingProps = {
  withContainer?: boolean
}

export function Loading({ withContainer = true }: LoadingProps) {
  if (withContainer) {
    return (
      <div className="container">
        <p>Loading...</p>
      </div>
    )
  }
  return <p>Loading...</p>
}
