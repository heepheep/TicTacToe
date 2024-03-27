import logo from './pages/logo.png';
import loading from './pages/pngwing.com.png'
import './App.css';
import newGameLogo from './pages/newGame.png'
import {BrowserRouter, Routes, Route, NavLink, useNavigate, Navigate} from 'react-router-dom'
import React, {useState} from 'react'
import axios from 'axios'
import {CookiesProvider, useCookies} from 'react-cookie'

let userData = {
  PlayerName: ""
}

function App() {
  const [cookies, setCookie] = useCookies(['user', 'gameid', 'logoicon']);
  return (
    <BrowserRouter>
      <div className="mainContainer">
        <header>
            <NavLink to={'/'}>
              { cookies.logoicon == "0" 
                ? <img src={logo} className="myStyle3"/> 
                : <img src={newGameLogo} className="myStyle3"/>
              }
              
            </NavLink>
        </header>
        <CookiesProvider>
          <Routes>
            <Route path="/" Component={LoginMenu}></Route>
            <Route path="/game" Component={Game}></Route>
            <Route path="/lobby" Component={Lobby}></Route>
          </Routes>
        </CookiesProvider>
        <footer>
            <div className="myStyle1">Разработчик: Черных Иван</div>
            <div className="myStyle2">Для ГБПОУ МО "Воскресенский колледж"</div>
        </footer>
      </div>
    </BrowserRouter>
    
  );
  function Lobby(){
    
    const navigate = useNavigate();
    if(cookies.gameid == -1){
      axios.post(`http://192.168.1.10:8080/make-game`, {
        playerName: cookies.user
      }).then((resp) => {
        if(resp.status == 200){
          setCookie('gameid', resp.data.id, {path: '/'});
          navigate('/game');
        } else if (resp.status == 201){
          setCookie('gameid', resp.data.id, {path: '/'});
        } else{
          console.log(resp.data)
        }
      })
    }
    React.useEffect(() => {
      async function getPost() {
        const response = await axios.get(`http://192.168.1.10:8080/get-game/${cookies.gameid}`).then((resp) => {
        if(resp.status == 200){
          navigate('/game');
        }
      })}
      const interval = setInterval(() => {
        getPost();
      }, 2000);
      return () => clearInterval(interval);
      
    }, []);
    return(
      <main className='containers'>
        <div className='loadContainer'>
          <img src={loading} className='load'></img>
        </div>
        <p className='loadText'>
          Идет поиск игры
        </p>
      </main>
    )
  }
  function LoginMenu(){
    setCookie('logoicon', "0", {path: '/'});
    const handleChange = event => {
      userData.PlayerName = event.target.value;
    }
    const handleClick = event => {
      setCookie('user', userData.PlayerName, {path: '/'});
      setCookie('gameid', -1, {path: '/'});
    }
    return(
      <main className="containers">
            <input id='nameSpace' className="nameInput myStyle4" placeholder="Имя:" onChange={handleChange}/>
            <NavLink to={'/lobby'}><input className="nameInput" value="Играть" type="button" onClick={handleClick}/></NavLink>
      </main>
    )
  }
  function Game(props){
    const [boarding, setBoard] = React.useState(null);
    const [mainData, setData] = React.useState(null);

    React.useEffect(() => {
      async function getPost() {
        const response = await axios.get(`http://192.168.1.10:8080/get-game/${cookies.gameid}`).then((resp) => {
        setData(resp.data);
        console.log(resp.data.ans)
        setBoard(resp.data.board);
        if(resp.data.isFinnished){
          setCookie('logoicon', "1", {path: '/'});
        }
      })}
      const interval = setInterval(() => {
        getPost();
      }, 1000);
      return () => clearInterval(interval);
      
    }, []);
    function renderSquare(index){
      if(mainData == null){
        return(
          <button className='cell'></button>
        )
      }
      else{
        function getSymbol(){
          if(boarding != undefined){
            return boarding[Math.floor(index/3)][index%3];
          } else{
            return "";
          }
          
        }
        async function updSymbol(e){
          await axios.post("http://192.168.1.10:8080/make-move",{id: cookies.gameid, player: cookies.user, x: (index%3), y: (Math.floor(index/3))}).then((resp) => {
          console.log(resp.data);
          if(resp.status == 230 || resp.status == 240){

          } else{
            console.log(resp.data)
            setBoard(resp.data.board.board)
          

          }
          
          });
        }
        
        return(
          <button className='cell' onClick={updSymbol}>{getSymbol()}</button>
        )
      }
    }
    
    return (
      <div className="containers">
        <div className="gameboard">
          {renderSquare(0)}
          {renderSquare(1)}
          {renderSquare(2)}
          {renderSquare(3)}
          {renderSquare(4)}
          {renderSquare(5)}
          {renderSquare(6)}
          {renderSquare(7)}
          {renderSquare(8)}
        </div>
        <div className="game-info">{
          
          mainData != null 
          ? <p className='game-info'>{mainData.ans}</p> 
          : console.log(mainData)}
        </div>
      </div>
    );
  };
}
export default App;
