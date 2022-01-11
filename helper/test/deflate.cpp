#include <stdio.h>
#include <stdint.h>
#include <vector>
#include <zlib.h>
#include <string.h>
#include <fstream>
#include <ctype.h>

using namespace std;

#define MAX_BLOCK_SIZE 8192 // the client will crash if you make this bigger, so don't.
#define WriteToBuffer(type, val, buffer, idx) if(idx + sizeof(type) > buffer.size()) { buffer.resize(idx + sizeof(type)); } *(type*)&buffer[idx] = val;  

uint32_t DeflateData(const char *buffer, uint32_t len, char *out_buffer, uint32_t out_len_max) {
	z_stream zstream;
	memset(&zstream, 0, sizeof(zstream));
	int zerror;

	zstream.next_in = const_cast<unsigned char*>(reinterpret_cast<const unsigned char*>(buffer));
	zstream.avail_in = len;
	zstream.zalloc = Z_NULL;
	zstream.zfree = Z_NULL;
	zstream.opaque = Z_NULL;
	deflateInit(&zstream, Z_FINISH);

	zstream.next_out = reinterpret_cast<unsigned char*>(out_buffer);
	zstream.avail_out = out_len_max;
	zerror = deflate(&zstream, Z_FINISH);

	if (zerror == Z_STREAM_END)
	{
		deflateEnd(&zstream);
		return (uint32_t)zstream.total_out;
	}
	else
	{
		zerror = deflateEnd(&zstream);
		return 0;
	}
}

bool Deflate(const std::vector<char> &file, std::vector<char> &out_buffer) {
	uint32_t pos = 0;
	uint32_t remain = (uint32_t)file.size();
	uint8_t block[MAX_BLOCK_SIZE + 128];
	while(remain > 0) {
		uint32_t sz;
		if (remain >= MAX_BLOCK_SIZE) {
			sz = MAX_BLOCK_SIZE;
			remain -= MAX_BLOCK_SIZE;
		} else {
			sz = remain;
			remain = 0;
		}

		uint32_t block_len = sz + 128;
		uint32_t deflate_size = (uint32_t)DeflateData(&file[pos], sz, (char*)&block[0], block_len);
		if(deflate_size == 0)
			return false;

		pos += sz;

		uint32_t idx = (uint32_t)out_buffer.size();
		WriteToBuffer(uint32_t, deflate_size, out_buffer, idx);
		WriteToBuffer(uint32_t, sz, out_buffer, idx + 4);
		out_buffer.insert(out_buffer.end(), block, block + deflate_size);
	}

	return true;
}


vector<char> readFile(const char* filename) {	
    vector<char> vec;
    ifstream file(filename, ios::binary);
	if (!file.is_open()) {
		printf("open %s failed: %s\n", filename, strerror(errno));
		return vec;
	}
    file.unsetf(ios::skipws);
    streampos fileSize;

    file.seekg(0, ios::end);
    fileSize = file.tellg();
    file.seekg(0, ios::beg);

    vec.reserve(fileSize);

    vec.insert(vec.begin(), istream_iterator<char>(file), istream_iterator<char>());
    return vec;
}

void hexdump(vector<char> buf) {
	int i, j;
	int sz = buf.size();
	printf("\n");   
	for (i=0; i<sz; i+=16) {
		printf("%08x  ", i);
		for (j=0; j<16; j++) {
			if (i+j < sz) {
				printf("%02x ", buf[i+j]);
			} else {
				printf("00 ");
			}
			if (j == 7) {
				printf(" ");
			}
		}
		printf("|");
		for (j=0; j<16; j++) {
			if (i+j < sz) {
				printf("%c", isprint(buf[i+j]) ? buf[i+j] : '.');
			}
		}
		printf("|\n");
	}
}

int main() {
    vector<char> data = readFile("test.txt");
    vector<char> out;
    if (!Deflate(data, out)) {
        printf("failed to deflate");
        return 1;
    }

	hexdump(out);
}

