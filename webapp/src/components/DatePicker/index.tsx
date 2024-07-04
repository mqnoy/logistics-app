import { FC, useEffect, useRef } from 'react';
import 'bulma/css/bulma.min.css';
import 'bulma-calendar/dist/css/bulma-calendar.min.css';
import bulmaCalendar from 'bulma-calendar';
import { DatePickerEventSelect } from './type';

type DatePickerProps = {
    onSelected: (e: DatePickerEventSelect) => void
    isRange?: boolean
}

export const DatePicker: FC<DatePickerProps> = ({ onSelected, isRange }) => {
    const calendarRef = useRef(null);

    useEffect(() => {
        if (calendarRef && calendarRef.current) {
            const calendarInstance = bulmaCalendar.attach(calendarRef.current, {
                isRange,
                type: 'date',
                dateFormat: 'MM-dd-yyyy',
                showButtons: false,
                showHeader: false,
                showFooter: false,
                closeOnSelect: true,
                labelFrom: 'from',
                labelTo: 'to'
            })[0];

            calendarInstance.on('select', date => {
                onSelected(date);
            });

            return () => {
                calendarRef.current = null
            };
        }
    }, []);

    return (
        <>
            <input
                className="input"
                type="date"
                ref={calendarRef}
            />
        </>
    );
};

