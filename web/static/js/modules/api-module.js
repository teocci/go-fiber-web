/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-Aug-18
 */
export default class APIModule {
    static gatherLimit(req) {
        if (isNil(req)) return

        const limit = isNil(req.limit) || req.limit === 0 ?
            isNil(pageInfo.limit) || pageInfo.limit === 0 ? null : pageInfo.limit : req.limit

        return isNil(limit) ? '' : `?limit=${limit}`
    }


    /** @typedef {Object} CommonRequest
     * @property {string} action
     * @property {string} uid
     * @property {string?} xsubject
     * @property {int?} limit
     */

    /**
     * @param req {CommonRequest}
     * @param callback {Function}
     */
    static requestFilters(req, callback) {
        if (isNil(req) || isNil(callback)) return

        const {action, uid} = req
        if (isNil(action) || isNil(uid)) return

        const limit = APIModule.gatherLimit(req)

        const url = `/api/v1/filters/${action}/${uid}${limit}`
        console.log('requestFilters', {url})

        APIModule.commonFetch(url).then(d => {
            if (isNil(d)) throw new Error('Filters not found')

            const {data} = d
            if (isNil(data)) throw new Error('Filters data not found')

            callback(data)
        }).catch(e => {
            console.error(e)
        })
    }

    /**
     * @param req {CommonRequest}
     * @param callback {Function}
     */
    static requestPositions(req, callback) {
        if (isNil(req) || isNil(callback)) return

        const {action, uid} = req
        if (isNil(action) || isNil(uid)) return

        const limit = APIModule.gatherLimit(req)
        const xsubject = isNil(req.xsubject) ? '' : `/${req.xsubject}`

        const url = `/api/v1/positions/${action}/${uid}${xsubject}${limit}`
        console.log('requestPositions', {url})

        APIModule.commonFetch(url).then(d => {
            if (isNil(d)) throw new Error('Null position response')

            const {data} = d
            if (isNil(data)) throw new Error('Positions data not found')

            callback(data)
        }).catch(e => {
            console.error(e)
        })
    }

    /** @typedef {Object} ProductRequest
     * @property {string} uid
     * @property {int} limit
     */

    /**
     * @param req {ProductRequest}
     * @param callback {Function}
     */
    static fetchProducts(req, callback) {
        if (isNil(req) || isNil(callback)) return

        const {uid, action} = req
        if (isNil(action) || isNil(uid)) return

        const limit = isNil(req.limit) ? '' : `?limit=${req.limit}`

        const url = `/api/v1/products/${action}/${uid}${limit}`
        console.log('fetchProducts', {url})

        APIModule.commonFetch(url).then(d => {
            if (isNil(d)) throw new Error('Null position response')

            const {data} = d
            if (isNil(data)) throw new Error('Positions data not found')

            callback(data)
        }).catch(e => {
            console.error(e)
        })

        // fetch(url)
        //     .then(res => {
        //         // no matching records found
        //         if (res.status === 404) return {data: []}
        //         if (res.ok) return res.json()
        //
        //         throw Error('oh no :(')
        //     })
        //     .then(d => {
        //         const data = d.data.products.map(product => [product.id, product.name, product.priceU / 100, product.salePriceU / 100, product])
        //
        //         this.data = []
        //         for (const product of d.data.products) {
        //             const {identical} = product
        //             if (isNil(identical)) continue
        //
        //             const item = {
        //                 id: product.id,
        //                 name: product.name,
        //                 prices: {
        //                     base: product.priceU / 100,
        //                     sale: product.salePriceU / 100,
        //                 },
        //                 competitors: [],
        //             }
        //
        //             for (const idem of identical) {
        //                 const competitor = {
        //                     id: idem.supplierInfo.id,
        //                     name: idem.supplierInfo.name,
        //                     base: idem.priceU / 100,
        //                     sale: idem.salePriceU / 100,
        //                 }
        //                 item.competitors.push(competitor)
        //             }
        //
        //             this.data.push(item)
        //         }
        //
        //         this.state = data.length > 0 ? STATE_DATA_LOADED : STATE_DATA_EMPTY
        //
        //         this.resolver(data)
        //     })
    }

    static async commonFetch(url) {
        const response = await fetch(url)
        if (!response.ok) return null
        const data = await response.json()
        if (isNil(data)) return null

        return data
    }
}