'use strict';

// const APIurl = 'http://backend:9000/';
const APIurl = 'http://localhost:9000/';
const fetchAPI = async (uri, headers) => {
    return fetch(`${APIurl}${uri}`, headers)
        .then(response => response.json())
        .then(data => {
            return data;
        });
};