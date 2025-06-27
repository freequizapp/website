import "./App.css";
import MultipleChoice from "../components/MultipleChoice";
import SearchBar from "../components/SearchBar";
function App() {
  return (
    <div className="home-background flex flex-col justify-center bg-black">
      <SearchBar />
      <MultipleChoice />
    </div>
  );
}

export default App;
