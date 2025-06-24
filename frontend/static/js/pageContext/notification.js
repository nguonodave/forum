export const notify = (message, color) => {
    const container = document.querySelector(".search-bar")

    let notification = document.createElement("div")
    notification.classList.add("notification");
    notification.style.backgroundColor = `${color}`;
    notification.innerHTML = `${message}`
    container.appendChild(notification);

    setTimeout(()=>{
        notification.style.display = "none";
    },3000);
}