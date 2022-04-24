'use strict';

const ListAllProducts = (props) => {
    const getData = async () => {
        fetchAPI('product', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            }
        })
            .then(response => {
                console.log(response);
                props.setText(JSON.stringify(response));
            })
    }
    return (
        <React.Fragment>
            <button onClick={getData}> List products</button>
        </React.Fragment>
    );
}
const GetProduct = (props) => {
    const getData = async () => {
        fetchAPI(`product/${props.ID}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            }
        })
            .then(response => {
                console.log(response);
                props.setText(JSON.stringify(response));
            })
    }
    return (
        <React.Fragment>
            <button onClick={getData}> List by ID</button>
        </React.Fragment>
    );
}
const DeleteProduct = (props) => {
    const update = async () => {
        fetchAPI(`product/${props.ID}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            }
        })
            .then(response => {
                console.log(response);
                props.setText(JSON.stringify(response));
            })
    }
    return (
        <React.Fragment>
            <button onClick={update}> Delete product with ID</button>
        </React.Fragment>
    );
}
const InsertProduct = (props) => {
    const update = async () => {
        let updatefields;
        if (!!props.Name) {
            updatefields = {
                ...updatefields,
                Name: props.Name
            }
        }
        if (!!props.Ammount) {
            updatefields = {
                ...updatefields,
                Ammount: props.Ammount
            }
        }
        let header = {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({ ...updatefields })
        }
        console.log(header);
        fetchAPI(`product`, header)
            .then(response => {
                console.log(response);
                props.setText(JSON.stringify(response));
            })
    }
    return (
        <React.Fragment>
            <button onClick={update}> Insert product</button>
        </React.Fragment>
    );
}

function testControls(props) {
    const [ResultAsText, setResultAsText] = React.useState("This result is empty");
    const [ID, setID] = React.useState("");
    const [Name, setName] = React.useState("");
    const [Ammount, setAmmount] = React.useState(0);
    return (
        <React.Fragment>
            <ListAllProducts setText={setResultAsText} />
            <br />
            ID: <input type="text" value={ID} onChange={(e) => setID(e.target.value)} />
            <br />
            <GetProduct ID={ID} setText={setResultAsText} />
            <br />
            <DeleteProduct ID={ID} setText={setResultAsText} />
            <br />
            Name: <input type="text" value={Name} onChange={(e) => setName(e.target.value)} />
            Ammount: <input type="number" value={Ammount} onChange={(e) => setAmmount(Number(e.target.value))} />
            <br />
            <InsertProduct Name={Name} Ammount={Ammount} setText={setResultAsText} />
            <br />
            <div>
                {ResultAsText}
            </div>

        </React.Fragment>
    );
}

const e = React.createElement;
const domContainer = document.querySelector('.container');
const root = ReactDOM.createRoot(domContainer);
root.render(e(testControls))
