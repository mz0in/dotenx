import React from "react"
import ReactDOM from "react-dom/client"
import "reactflow/dist/style.css"
import { App } from "./features/app"
import "./index.css"

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
	<React.StrictMode>
		<App />
	</React.StrictMode>
)
