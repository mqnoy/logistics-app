import { FC } from "react"
import Layout from "../Layout"
import { GoodsList } from "../../components"

export const GoodsPage: FC = () => {
    return (
        <Layout>
            <GoodsList />
        </Layout>
    )
}

export default GoodsPage