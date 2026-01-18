export function getParams(): { tags: string[]; recipe: string } {
  const params = new URLSearchParams(window.location.search)
  const tagParam = params.get('tags') ?? ''
  return {
    tags: tagParam ? tagParam.split(',') : [],
    recipe: params.get('recipe') ?? '',
  }
}

export function setParams(tags: string[], recipe: string) {
  const params = new URLSearchParams()
  if (tags.length > 0) params.set('tags', tags.join(','))
  if (recipe !== '') params.set('recipe', recipe)
  const search = params.toString()
  const url = search === '' ? '/browse' : `/browse?${search}`
  window.history.replaceState(null, '', url)
}
