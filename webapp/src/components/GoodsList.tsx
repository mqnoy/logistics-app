
import { FC, ReactNode, useEffect, useState } from "react";
import { Goods, ModalActionGoods } from "../types";
import { TableCustom } from "./TableCustom";
import { useLazyGetDetailGoodQuery, useLazyGetListGoodsQuery } from "../api";
import { rtkUtils, toastUtils } from "../utils";
import { FaPlus } from "react-icons/fa6";
import { Modal } from "./Modal";
import { GoodsForm } from ".";

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

    const [isModalActiveCU, setIsModalActiveCU] = useState(false);
    const [modalTitle, setModalTitle] = useState('')
    const [action, setAction] = useState<ModalActionGoods>();

    const showModalCU = () => {
        setIsModalActiveCU(true);
    };

    const closeModal = () => {
        setIsModalActiveCU(false);
    };


    const [getDetail, { data: dataGetDetail, error: errorGetDetail, isLoading: isloadingGetDetail }] = useLazyGetDetailGoodQuery();
    const handleActionUpdate = (props: Goods) => {
        getDetail(props.id)
        setModalTitle("Edit Goods")
        setAction("update")
        showModalCU()
    }

    useEffect(() => {
        if (errorGetDetail) {
            const errorApi = rtkUtils.parseErrorRtk(errorGetDetail);
            toastUtils.fireToastError(errorApi)
        } else if (isloadingGetDetail) {
            console.debug('loading..');
        }
    }, [errorGetDetail, isloadingGetDetail])


    return (
        <div className="section">
            <Modal
                title={modalTitle}
                isActive={isModalActiveCU}
                onClose={closeModal}
                content={
                    <GoodsForm
                        action={action}
                        dataDetail={dataGetDetail?.data}
                        actionIsDone={(isDone, error) => {
                            if (isDone && !error) {
                                toastUtils.fireToastSuccess("successfully", {
                                    onClose() {
                                        setIsModalActiveCU(false)
                                        trigger(Object.assign({
                                            limit: 10,
                                            offset: 0,
                                            page: page,
                                            orders: 'id desc',
                                        }));
                                    },
                                })
                            }

                            if (error) {
                                toastUtils.fireToastError(error)
                            }
                        }}
                    />
                }
            />
            <div className="columns">
                <div className="column">
                    <span className="icon-text has-text-info">
                        <h5 className="title">Goods</h5>
                        <button className="button is-primary is-small"
                            aria-label="add new goods"
                            onClick={() => {
                                setModalTitle("Add new Goods")
                                setAction("create")
                                showModalCU()
                            }}
                        > <FaPlus className="has-text-white" /></button>
                    </span>

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
                                                handleActionUpdate(item)
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
