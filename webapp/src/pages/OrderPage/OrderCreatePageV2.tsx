import { FC, useEffect, useState } from 'react'
import { Controller, SubmitHandler, useFieldArray, useForm } from 'react-hook-form'
import Layout from '../Layout'
import { MultipleOrderCreateRequest, OrderItemRequest, OrderTypeEnum } from '../../types'
import { useLazyGetListGoodsQuery, usePostMultipleOrderCreateMutation } from '../../api'
import { rtkUtils, toastUtils } from '../../utils'
import { useDebounce } from 'use-debounce'
import { FaTrash } from 'react-icons/fa6'
import Select from 'react-select'

type IFormOrder = {
    items: OrderItemRequest[]
    type: number
}

type Option = {
    value: number
    label: string
}

export const OrderCreatePageV2: FC = () => {
    const {
        register,
        control,
        reset,
        formState: { errors },
        handleSubmit,
        setValue,
    } = useForm<IFormOrder>({
        defaultValues: {
            items: [
                {
                    code: '',
                    total: 0,
                },
            ],
        },
        criteriaMode: 'all',
    })
    const { fields, append, remove } = useFieldArray({
        control,
        name: 'items',
    })

    const [dataGoods, setDataGoods] = useState<Option[]>()
    const [keyword, setKeyword] = useState('')
    const [keywordDebounce] = useDebounce(keyword, 800)
    const [getListGoods, { data: dataListGoods, error: errorGoods, isLoading: isLoadingGoods }] =
        useLazyGetListGoodsQuery()

    useEffect(() => {
        if (dataListGoods) {
            const options = dataListGoods.data.rows.map((item, index) => ({
                value: index,
                label: item.code,
            }))
            setDataGoods(options)
        } else if (errorGoods) {
            const errorApi = rtkUtils.parseErrorRtk(errorGoods)
            toastUtils.fireToastError(errorApi)
        } else if (isLoadingGoods) {
            console.log(isLoadingGoods)
        }
    }, [dataListGoods, errorGoods, isLoadingGoods])

    useEffect(() => {
        getListGoods(
            Object.assign(
                {
                    limit: 10,
                    offset: 0,
                    page: 1,
                    orders: 'id desc',
                },
                {
                    keyword: keywordDebounce,
                }
            )
        )
    }, [keywordDebounce])

    const [
        postMultipleOrderCreate,
        {
            isLoading: isLoadingOrderCreate,
            isSuccess: isSuccessOrderCreate,
            error: errorOrderCreate,
        },
    ] = usePostMultipleOrderCreateMutation()
    useEffect(() => {
        if (isSuccessOrderCreate) {
            toastUtils.fireToastSuccess('successfully', {
                onClose() {
                    handleReset()
                },
            })
        } else if (errorOrderCreate) {
            const errorApi = rtkUtils.parseErrorRtk(errorOrderCreate)
            toastUtils.fireToastError(errorApi)
        } else if (isLoadingOrderCreate) {
            console.debug('loading...')
        }
    }, [isSuccessOrderCreate, errorOrderCreate, isLoadingOrderCreate])

    const onsubmit: SubmitHandler<IFormOrder> = (data: IFormOrder) => {
        console.debug('onsubmited', data)

        const orderItem = data.items.reduce<OrderItemRequest[]>(
            (result: OrderItemRequest[], curr: OrderItemRequest) => {
                result.push({
                    code: curr.code,
                    total: Number(curr.total),
                })
                return result
            },
            []
        )
        const payload: MultipleOrderCreateRequest = {
            items: orderItem,
            type: Number(data.type),
        }
        postMultipleOrderCreate(payload)
    }

    const handleReset = () => {
        reset()
    }

    const handleOnChange = (index: number, item: Option) => {
        setValue(`items.${index}.code`, item.label)
    }

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
                <div className="columns">
                    <div className="column">
                        <button
                            className="button is-primary has-text-white"
                            type="button"
                            onClick={() => {
                                append({ code: '', total: 0 })
                            }}
                        >
                            Add Item
                        </button>
                    </div>
                </div>
                <form className="" onSubmit={handleSubmit(onsubmit)}>
                    <div className="field">
                        <label htmlFor="" className="label">
                            Type
                        </label>
                        <div className="control">
                            <div className="select">
                                <select
                                    {...register('type', { min: 1 })}
                                    onChange={(e) => {
                                        setValue('type', Number(e.target.value))
                                    }}
                                >
                                    <option value={0}>- select type -</option>
                                    <option value={OrderTypeEnum.ORDER_IN}>In</option>
                                    <option value={OrderTypeEnum.ORDER_OUT}>Out</option>
                                </select>
                            </div>
                        </div>
                        {errors.type?.type === 'min' && (
                            <p className="help is-danger">select order type</p>
                        )}
                    </div>
                    {fields.map((field, index) => (
                        <div key={field.id} className="columns">
                            <div className="column is-1">
                                <h5>Item{`${index}`}</h5>
                                <button
                                    className="button is-danger"
                                    type="button"
                                    onClick={() => remove(index)}
                                >
                                    <FaTrash />
                                </button>
                            </div>
                            <div className="column is-11">
                                <div className="field">
                                    <label htmlFor="" className="label">
                                        Code
                                    </label>
                                    <div className="control">
                                        <Controller
                                            name={`items.${index}.code`}
                                            control={control}
                                            render={() => {
                                                return (
                                                    <Select
                                                        required={true}
                                                        isClearable={true}
                                                        isSearchable={true}
                                                        options={dataGoods}
                                                        onInputChange={(newValue: string) => {
                                                            setKeyword(newValue)
                                                        }}
                                                        onChange={(item) => {
                                                            if (item) {
                                                                handleOnChange(index, item)
                                                            }
                                                        }}
                                                    />
                                                )
                                            }}
                                        />
                                    </div>
                                </div>

                                <div className="field">
                                    <label htmlFor={`items.${index}.total`} className="label">
                                        Total
                                    </label>
                                    <div className="control">
                                        <input
                                            {...register(`items.${index}.total`, { min: 1 })}
                                            type="number"
                                            placeholder="0"
                                            className={`input ${errors.items?.[index]?.total && 'is-danger'}`}
                                            required
                                        />
                                    </div>
                                    {errors.items?.[index]?.total?.type === 'min' && (
                                        <p className="help is-danger">total should minimum 1</p>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))}

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
