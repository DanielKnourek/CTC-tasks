'use strict';

// const APIurl = 'http://backend:9000/';
const APIurl = 'api/';
// const APIurl = 'http://localhost:8080/api/';
const fetchAPI = async (uri, headers) => {
    return fetch(`${APIurl}${uri}`, headers)
        .then(response => response.json())
        .then(data => {
            return data;
        });
};