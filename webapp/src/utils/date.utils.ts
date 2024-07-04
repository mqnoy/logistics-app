import { DateTime } from 'luxon'

export const epochToDateTime = (e: number): DateTime => {
    return DateTime.fromSeconds(e)
}

export const getCurrentTz = (): string => {
    const { timeZone } = Intl.DateTimeFormat().resolvedOptions()
    return timeZone
}

export const epochNumberToDateTimeStr = (e: number): string => {
    const timeZone = getCurrentTz()
    const dt = DateTime.fromSeconds(e)
    dt.setZone(timeZone)
    return dt.toFormat('MM-dd-yyyy HH:mm:ss')
}

export const dateToEpoch = (d: Date): number => {
    return Math.floor(DateTime.fromJSDate(d).toMillis() / 1000)
}
