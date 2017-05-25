((window) => {
    function unlinkFile(path) {
        if (!confirm('Are you sure you want to delete this file?')) {
            return;
        }
        fetch(path, {
                method: "DELETE"
            })
            .then(
                (response) => {
                    console.log(response.json());
                },
                (error) => {
                    console.error(error.message);
                    alert(error.message);
                });
    }

    window.unlinkFile = unlinkFile;
})(window);