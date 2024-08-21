export const convertDateToInput = (dateObject: number) => {
  if (dateObject === -1) return
  const date = new Date(dateObject)
  const year = date.getFullYear()

  let month: number | string = date.getMonth() + 1
  let day: number | string = date.getDate()

  if (month < 10) month = '0' + month
  if (day < 10) day = '0' + day

  return `${year}-${month}-${day}`
}
