import axios from 'axios';

class BiddrClient{

    userAuth = '';
    domain = '';

    constructor(baseUrl){
        this.domain = baseUrl
    }
    
    async requestConfig(url, method, {headers, body} = {}){
        const myParameters = {
            url : url, 
            method : method, 
            baseURL : this.domain,  
            validateStatus :  status => status < 500,
            headers: {
                ...(headers || {}),
                authorization: `Bearer ${this.userAuth}`,
            }
        };

        if(body){
            myParameters['data'] = body;
        };
        return await axios(myParameters);
    }

    async createUser(user, display, pass){
        const response = await this.requestConfig('/api/v1/users/', 'post', {
            body: {
            username: user,
            password: pass,
            displayName: display
        }});

        return response.data;
    }

    async getUserInformation(userName){
        const encodedUser = encodeURIComponent(userName)
        const response = await  this.requestConfig(`/api/v1/users/${encodedUser}`, 'get');
        return response.data;
    }

    async userLogIn(username, password){
        const response = await this.requestConfig('/api/v1/users/login', 'post',{
            body: {
            username: username,
            password: password,
            
        }});
        const data = response.data;
        this.userAuth = data['authToken'];
        return data;
    }

    async userLogout() {
        this.authToken = '';
    }

    async getHighestBidForAll(){
        const response = await this.requestConfig('/api/v1/auctions/bids', 'get');
        return response.data;
     }

    async getHighestBidForOneItem(name){
        const encodedItem = encodeURIComponent(name)
        const response = await this.requestConfig(`/api/v1/auctions/bids/${encodedItem}`, 'get');
        return response.data;
    }

    async makeBid(name, bid){
        const encodedItem = encodeURIComponent(name)
        const response = await this.requestConfig(`/api/v1/auctions/bids/${encodedItem}`, 'post', {
            body: {bidAmount: bid},
        });


        if (response.status === 403) {
            throw new Error('user does not have permissions to place a bid');
        } else if (response.status === 400) {
            throw new Error(response.data.message);
        }
        return response.data;
    }
    
    async getAllItems(){
        const response = await this.requestConfig(`/api/v1/auctions/items`, 'get');
        return response.data;
    }

    async createNewItem(description, image, name){
        console.log(name);
        const response = await this.requestConfig('/api/v1/auctions/items', 'post', {
            body:
            {   description,
                image,
                name
            },
        });
        return response.data;
    }

    async deleteItem(name){
        const encodedItem = encodeURIComponent(name)
        const response = await this.requestConfig(`/api/v1/auctions/items/${encodedItem}`, 'delete');
        return response.data;
    }

    async getSpecificItem(name){
        const encodedItem = encodeURIComponent(name)
        const response = await this.requestConfig(`/api/v1/auctions/items/${encodedItem}`, 'get');
        return response.data;
    }

    async updateItem(description, image, name){
        const encodedItem = encodeURIComponent(name)
        const response = await this.requestConfig(`/api/v1/auctions/items/${encodedItem}`, 'put', {
            body: {
                description: description,
                image: image,
                name: name
            }
        });
        return response.data;
    }

    async establishWebSocket(){
        const response = await this.requestConfig('/api/v1/ws', 'get');
        return response.data;
    }

}

export default new BiddrClient("http://localhost:8080")
