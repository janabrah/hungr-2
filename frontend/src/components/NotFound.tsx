type NotFoundProps = {
  title?: string
  message?: string
}

export function NotFound({
  title = 'Not Found',
  message = 'The page you are looking for does not exist.',
}: NotFoundProps) {
  return (
    <div className="container">
      <h1>{title}</h1>
      <p>{message}</p>
    </div>
  )
}
