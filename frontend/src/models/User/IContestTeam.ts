export interface IContestTeam {
  id: number
  owner_id: number
  contest_id: number
  name: string
  description: string
  approved?: number
}

export interface ITeamParticipant {
  id: number
  login: string
  name: string
  prevPoints: number
  tg: string
  role: string
  gpa: number
}
