import { DateTime } from "luxon"

export const epochToDateTime = (e: number): DateTime => {
    return DateTime.fromSeconds(e)
}