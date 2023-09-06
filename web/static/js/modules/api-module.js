/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-Aug-18
 */
export default class APIModule {
    /** @typedef {Object} CommonRequest
     * @property {string} action
     * @property {string} sellerId
     * @property {string?} categoryId
     * @property {string?} xsubject
     * @property {int?} limit
     */

    /** @typedef {Object} ProductRequest
     * @property {string} uid
     * @property {int?} limit
     */

    /** @typedef {Object} SellerRequest
     * @property {string} sellerId
     */

    /** @typedef {Object} AdsListRequest
     * @property {string} sellerId
     * @property {int?} limit
     */

    /** @typedef {Object} MarketingControlRequest
     * @property {string} sellerId
     * @property {string} action
     * @property {string} campaignId
     */

    static gatherLimit(req) {
        if (isNil(req)) return

        const limit = isNil(req.limit) || req.limit === 0 ?
            isNil(pageInfo.limit) || pageInfo.limit === 0 ? null : pageInfo.limit : req.limit

        return isNil(limit) ? '' : `?limit=${limit}`
    }

    /**
     *
     * @param req {MarketingControlRequest}
     * @param callback
     */
    static requestMarketingControl(req, callback) {
        if (isNil(req) || isNil(callback)) throw new Error('Null request or callback')

        const {sellerId} = req || pageInfo.sellerId
        if (isNil(sellerId)) throw new Error('Invalid seller id: null')

        const {action, campaignId} = req
        if (isNil(action)) throw new Error('Invalid action: null')
        if (isNil(campaignId)) throw new Error('Invalid campaign id: null')

        const url = `/api/v1/marketing/${sellerId}/${action}/${campaignId}`
        console.log('requestMarketingControl', {url})

        APIModule.commonFetch(url).then(d => {
            if (isNil(d)) throw new Error('Invalid response: null')
            console.log('requestMarketingControl', {d})

            callback(d)
        }).catch(e => {
            console.error(e)
        })
    }

    /**
     *
     * @param req {AdsListRequest}
     * @param callback
     */
    static requestCampaignList(req, callback) {
        if (isNil(req) || isNil(callback)) throw new Error('Null request or callback')

        const {sellerId} = req || pageInfo.sellerId
        if (isNil(sellerId)) throw new Error('Invalid seller id: null')

        const limit = APIModule.gatherLimit(req)

        const url = `/api/v1/marketing/${sellerId}${limit}`
        console.log('requestAdsList', {url})

        APIModule.commonFetch(url).then(d => {
            if (isNil(d)) throw new Error('Ads list not found')

            const {data} = d
            if (isNil(data)) throw new Error('Ads list data not found')

            callback(data)
        }).catch(e => {
            console.error(e)
        })
    }

    /**
     *
     * @param req {SellerRequest}
     * @param callback
     */
    static requestSeller(req, callback) {
        if (isNil(req) || isNil(callback)) throw new Error('Null request or callback')

        const {sellerId} = req || pageInfo.sellerId
        if (isNil(sellerId)) throw new Error('Invalid seller id: null')

        const url = `/api/v1/seller/${sellerId}`
        console.log('requestSeller', {url})

        APIModule.commonFetch(url).then(d => {
            if (isNil(d)) throw new Error('Seller\'s details not found')

            const {data} = d
            if (isNil(data)) throw new Error('Seller\'s details data not found')

            callback(data)
        }).catch(e => {
            console.error(e)
        })
    }

    /**
     * @param req {CommonRequest}
     * @param callback {Function}
     */
    static requestFilters(req, callback) {
        if (isNil(req) || isNil(callback)) throw new Error('Null request or callback')

        const {action} = req
        if (isNil(action)) throw new Error('Invalid action: null')

        const uid = action === 'category' ? req.categoryId : req.sellerId
        if (isNil(uid)) throw new Error('Null uid')

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
        if (isNil(req) || isNil(callback)) throw new Error('Null request or callback')

        const {action, sellerId} = req
        if (isNil(action) || isNil(sellerId)) throw new Error('Action or seller id not found')

        let cat = ''
        if (req.action === 'category') {
            if (isNil(req.categoryId)) throw new Error('Category id not found')
            cat = `/cat/${req.categoryId}`
        }

        const xsubject = isNil(req.xsubject) ? '' : `/${req.xsubject}`

        const limit = APIModule.gatherLimit(req)

        const url = `/api/v1/positions/${action}/${sellerId}${cat}${xsubject}${limit}`
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