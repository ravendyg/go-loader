const fs = require('fs');
const path = require('path');
const http = require('http');

const server = http.createServer((req, res) => {
    const {
        url,
        headers: {
            range,
        },
    } = req;
    let rangeStart = 0;
    let rangeEnd;
    if (range) {
        [rangeStart, rangeEnd] = range.split('-');
    }
    let fileName = path.join(__dirname, url);
    if (fileName[fileName.length - 1] === path.sep) {
        fileName += 'main.go';
    }
    fs.exists(fileName, exists => {
        if (!exists) {
            res.statusCode = 404;
            return res.end();
        }

        fs.stat(fileName, (statErr, stat) => {
            if (statErr) {
                console.error(statErr);
                res.statusCode = 500;
                return res.end();
            }

            fs.open(fileName, 'r', (openErr, fd) => {
                if (openErr) {
                    console.error(openErr);
                    res.statusCode = 500;
                    return res.end();
                }

                if (!rangeEnd || stat.size < rangeEnd) {
                    rangeEnd = stat.size;
                }
                const chunkSize = rangeEnd - rangeStart;
                const buffer = new Buffer(chunkSize);
                fs.read(fd, buffer, rangeStart, chunkSize, 0, (readErr, read) => {
                    if (readErr) {
                        console.error(readErr);
                        res.statusCode = 500;
                        return res.end();
                    }

                    res.write(buffer, writeError => {
                        if (writeError) {
                            console.error(writeError);
                            res.statusCode = 500;
                            return res.end();
                        }

                        console.log(`sent ${read} bytes`);
                        res.end();
                    });
                });
            });
        });
    });
});

server.listen(3011);
