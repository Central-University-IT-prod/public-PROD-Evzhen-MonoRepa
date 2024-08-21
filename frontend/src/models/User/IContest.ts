export interface IContest {
  id: number
  name: string
  description: string
  field: string
  start_date: number
  end_date: number
  min_teammates: number
  max_teammates: number
  profiles: []
  commands: []
  admin_id: number
  end: boolean
}

export interface IContestRequest {
  name: string
  description: string
  field: string
  start_date: number
  end_date: number
}
