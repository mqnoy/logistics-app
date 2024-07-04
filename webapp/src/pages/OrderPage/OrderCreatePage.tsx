import { FC, useEffect, useState } from 'react'
import { SubmitHandler, useForm } from 'react-hook-form'
import { DropdownSearch } from '../../components/DropdownSearch'
import Layout from '../Layout'
import { BaseResponse, Goods, ListResponse, OrderTypeEnum } from '../../types'
import mockGoods from '@assets/mock/goods.json'

type IFormOrder = {
    goodsCode: string,
    total: number
    type: number
}

export const OrderCreatePage: FC = () => {
    const { register, formState: { errors }, handleSubmit, watch, setValue } = useForm<IFormOrder>({
        defaultValues: {
            goodsCode: '',
            total: 0,
            type: 0,
        },
        criteriaMode: "all"
    })

    const total = watch('total')

    const onsubmit: SubmitHandler<IFormOrder> = (data: IFormOrder) => {
        console.debug('onsubmited', data)
    }

    // TODO: Implement service orderApi
    const raw = mockGoods as BaseResponse<ListResponse<Goods>>
    console.log(raw.data);
    const [mockDataGoods, setMockDataGoods] = useState<Goods[]>([])

    useEffect(() => {
        setMockDataGoods(raw.data.rows)
    }, [])


    return (
        <Layout>
            <div className="pt-5">
                <div className="columns">
                    <div className="column">
                        <span className="icon-text has-text-info">
                            <h5 className="title">Create Order</h5>
                        </span>

                    </div>
                </div>
                <form className="" onSubmit={handleSubmit(onsubmit)}>
                    <div className="field">
                        <label htmlFor="" className="label">
                            Code
                        </label>
                        <div className="control">
                            <DropdownSearch
                                isRequired={true}
                                name="goodsCode"
                                items={mockDataGoods}
                                renderItem={(item: Goods) => {
                                    return item.code
                                }}
                                onSelected={(selectedData: Goods) => {
                                    console.log(selectedData);
                                    setValue('goodsCode', selectedData.code)
                                }}
                                onSearch={(v: string) => {
                                    // TODO: invoke goodsApi lazy query
                                    console.log(v);
                                }}
                            />

                        </div>
                        {errors.goodsCode?.type === "required" && (
                            <p className="help is-danger">code is required</p>
                        )}
                    </div>
                    <div className="field">
                        <label htmlFor="" className="label">
                            Total
                        </label>
                        <div className="control">
                            <input
                                {...register("total", { min: 1 })}
                                type="number"
                                value={total}
                                onChange={(e) => {
                                    setValue('total', Number(e.target.value))
                                }}
                                placeholder="0"
                                className={`input ${errors.total && 'is-danger'}`}
                                required
                            />
                        </div>
                        {errors.total?.type === "min" && (
                            <p className="help is-danger">total should minimum 1</p>
                        )}
                    </div>
                    <div className="field">
                        <label htmlFor="" className="label">
                            Type
                        </label>
                        <div className="control">
                            <div className="select">
                                <select
                                    {...register("type", { min: 1 })}
                                    onChange={(e) => {
                                        setValue('type', Number(e.target.value))
                                    }}
                                >
                                    <option value={0} >- select type -</option>
                                    <option value={OrderTypeEnum.ORDER_IN}>In</option>
                                    <option value={OrderTypeEnum.ORDER_OUT}>Out</option>
                                </select>
                            </div>
                        </div>
                        {errors.type?.type === "min" && (
                            <p className="help is-danger">select order type</p>
                        )}
                    </div>
                    <div className="field">
                        <button
                            className="button is-success is-fullwidth has-text-white"
                            type="submit"
                        >
                            Submit
                        </button>
                    </div>
                </form>
            </div>
        </Layout>

    )
}
