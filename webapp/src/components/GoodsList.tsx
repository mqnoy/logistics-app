
import { FC, ReactNode, useEffect, useState } from "react";
import { Goods } from "../types";
import { TableCustom } from "./TableCustom";
import { useLazyGetListGoodsQuery } from "../api";
import { rtkUtils, toastUtils } from "../utils";

type GoodsListProps = unknown

export const GoodsList: FC<GoodsListProps> = () => {
    const [keyword, setKeyword] = useState('')
    const [page, setPage] = useState(1);
    const [trigger, { data: goodsData, error, isLoading }] = useLazyGetListGoodsQuery();

    if (error) {
        const errorApi = rtkUtils.parseErrorRtk(error);
        toastUtils.fireToastError(errorApi)
    }

    if (isLoading) {
        console.log(isLoading);
    }

    const handleOnPageChange = (page: number) => {
        setPage(page)
    }

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setKeyword(event.target.value);
    };

    const handleSearch = () => {
        setPage(1);
        trigger(Object.assign({
            limit: 10,
            offset: 0,
            page: page,
            orders: 'id desc',
        }, {
            keyword: keyword
        },));
    }

    useEffect(() => {
        trigger(Object.assign({
            limit: 10,
            offset: 0,
            page: page,
            orders: 'id desc',
        }));
    }, [page])


    useEffect(() => {
        trigger(Object.assign({
            limit: 10,
            offset: 0,
            page: page,
            orders: 'id desc',
        }));
    }, [])

    const renderIsActive = (data: boolean): ReactNode => {
        if (data) {
            return <span className="button is-info is-small has-text-white">Active</span>
        }

        return <span className="button is-danger is-small has-text-white">Inactive</span>
    }

    return (
        <div className="section">
            <div className="columns">
                <div className="column">
                    <h5 className="title">Goods</h5>
                </div>
                <div className="column is-flex is-justify-content-flex-end">
                    <div className="field has-addons">
                        <div className="control">
                            <input
                                className="input"
                                type="text"
                                placeholder="Search by name or code"
                                value={keyword}
                                onChange={handleInputChange}
                            />
                        </div>
                        <div className="control">
                            <button
                                className="button is-primary has-text-white"
                                onClick={handleSearch}
                            >
                                Search
                            </button>
                        </div>
                    </div>
                </div>
            </div>
            {goodsData?.data &&
                <TableCustom
                    key={"a"}
                    onPageChange={handleOnPageChange}
                    data={goodsData.data}
                    tableHead={
                        <>
                            <th>code</th>
                            <th>name</th>
                            <th>active</th>
                            <th>stock</th>
                            <th>Action</th>
                        </>
                    }
                    renderRow={(item: Goods) => {
                        return < tr key={item.id} >
                            <td>{item.code}</td>
                            <td>{item.name}</td>
                            <td>{renderIsActive(item.is_active)}</td>
                            <td>{item.stock.total}</td>
                            <td>
                                <div className="field is-grouped">
                                    <p className="control">
                                        <button
                                            className="button is-primary is-outlined"
                                            onClick={() => {

                                            }} >
                                            view
                                        </button>
                                    </p>
                                    <p className="control">
                                        <button
                                            className="button is-primary is-outlined"
                                            onClick={() => {

                                            }} >
                                            edit
                                        </button>
                                    </p>
                                    <p className="control">
                                        <button
                                            className="button is-primary is-outlined"
                                            onClick={() => {

                                            }} >
                                            delete
                                        </button>
                                    </p>
                                </div>
                            </td>
                        </tr>
                    }}
                />
            }
        </div >
    )
}
