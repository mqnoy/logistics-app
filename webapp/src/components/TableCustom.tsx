import { FC, ReactNode, useState } from "react";
import { ListResponse } from "../types";

interface TableProps<T> {
    data: ListResponse<T>;
    tableHead: ReactNode;
    renderRow: (item: T, index: number) => ReactNode;
    onPageChange: (page: number) => void
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const TableCustom: FC<TableProps<any>> = ({ data, tableHead, renderRow, onPageChange }) => {
    const { metadata, rows } = data;

    const [currentPage, setCurrentPage] = useState(1);

    const handlePageChange = (page: number) => {
        setCurrentPage(page);
        onPageChange(page);
    };

    return (
        <div className="table-container">
            <table className="table is-fullwidth is-hoverable">
                <thead>
                    <tr className="is-uppercase">
                        {tableHead}
                    </tr>
                </thead>
                <tbody>
                    {rows.map((item, index) => (
                        renderRow(item, index)
                    ))}
                </tbody>
            </table>
            {rows.length === 0 && (
                <p className="has-text-centered mt-4">No data available.</p>
            )}
            <nav className="pagination" role="navigation" aria-label="pagination">
                <ul className="pagination-list">
                    {Array.from({ length: metadata.total_pages || 0 }, (_, i) => (
                        <li key={i}>
                            <button
                                key={i}
                                className={`pagination-link ${currentPage === i + 1 ? 'is-current' : ''}`}
                                aria-label={`page ${i + 1}`}
                                onClick={() => {
                                    handlePageChange(i + 1)
                                }}>
                                {i + 1}
                            </button>
                        </li>
                    ))}
                </ul>
            </nav>
        </div>
    );
}