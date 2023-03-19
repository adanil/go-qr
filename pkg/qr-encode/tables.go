// nolint:dupl
package qr_encode

var (
	versionSize = map[CodeLevel][40]int{
		L: {152, 272, 440, 640, 864, 1088, 1248, 1552, 1856, 2192,
			2592, 2960, 3424, 3688, 4184, 4712, 5176, 5768, 6360, 6888,
			7456, 8048, 8752, 9392, 10208, 10960, 11744, 12248, 13048, 13880,
			14744, 15640, 16568, 17528, 18448, 19472, 20528, 21616, 22496, 23648},
		M: {128, 224, 352, 512, 688, 864, 992, 1232, 1456, 1728,
			2032, 2320, 2672, 2920, 3320, 3624, 4056, 4504, 5016, 5352,
			5712, 6256, 6880, 7312, 8000, 8496, 9024, 9544, 10136, 10984,
			11640, 12328, 13048, 13800, 14496, 15312, 15936, 16816, 17728, 18672},
		Q: {104, 176, 272, 384, 496, 608, 704, 880, 1056, 1232,
			1440, 1648, 1952, 2088, 2360, 2600, 2936, 3176, 3560, 3880,
			4096, 4544, 4912, 5312, 5744, 6032, 6464, 6968, 7288, 7880,
			8264, 8920, 9368, 9848, 10288, 10832, 11408, 12016, 12656, 13328},
		H: {72, 128, 208, 288, 368, 480, 528, 688, 800, 976,
			1120, 1264, 1440, 1576, 1784, 2024, 2264, 2504, 2728, 3080,
			3248, 3536, 3712, 4112, 4304, 4768, 5024, 5288, 5608, 5960,
			6344, 6760, 7208, 7688, 7888, 8432, 8768, 9136, 9776, 10208},
	}

	numberOfBlocks = map[CodeLevel][40]int{
		L: {1, 1, 1, 1, 1, 2, 2, 2, 2, 4,
			4, 4, 4, 4, 6, 6, 6, 6, 7, 8,
			8, 9, 9, 10, 12, 12, 12, 13, 14, 15,
			16, 17, 18, 19, 19, 20, 21, 22, 24, 25},
		M: {1, 1, 1, 2, 2, 4, 4, 4, 5, 5,
			5, 8, 9, 9, 10, 10, 11, 13, 14, 16,
			17, 17, 18, 20, 21, 23, 25, 26, 28, 29,
			31, 33, 35, 37, 38, 40, 43, 45, 47, 49},
		Q: {1, 1, 2, 2, 4, 4, 6, 6, 8, 8,
			8, 10, 12, 16, 12, 17, 16, 18, 21, 20,
			23, 23, 25, 27, 29, 34, 34, 35, 38, 40,
			43, 45, 48, 51, 53, 56, 59, 62, 65, 68},
		H: {1, 1, 2, 4, 4, 4, 5, 6, 8, 8,
			11, 11, 16, 16, 18, 16, 19, 21, 25, 25,
			25, 34, 30, 32, 35, 37, 40, 42, 45, 48,
			51, 54, 57, 60, 63, 66, 70, 74, 77, 81},
	}

	numberOfCorrectionBytes = map[CodeLevel][40]int{
		L: {7, 10, 15, 20, 26, 18, 20, 24, 30, 18,
			20, 24, 26, 30, 22, 24, 28, 30, 28, 28,
			28, 28, 30, 30, 26, 28, 30, 30, 30, 30,
			30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
		M: {10, 16, 26, 18, 24, 16, 18, 22, 22, 26,
			30, 22, 22, 24, 24, 28, 28, 26, 26, 26,
			26, 28, 28, 28, 28, 28, 28, 28, 28, 28,
			28, 28, 28, 28, 28, 28, 28, 28, 28, 28},
		Q: {13, 22, 18, 26, 18, 24, 18, 22, 20, 24,
			28, 26, 24, 20, 30, 24, 28, 28, 26, 30,
			28, 30, 30, 30, 30, 28, 30, 30, 30, 30,
			30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
		H: {17, 28, 22, 16, 22, 28, 26, 26, 24, 28,
			24, 28, 22, 24, 24, 30, 28, 28, 26, 28,
			30, 24, 30, 30, 30, 30, 30, 30, 30, 30,
			30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
	}

	polynomialCoefficients = map[int][]int{
		7:  {87, 229, 146, 149, 238, 102, 21},
		10: {251, 67, 46, 61, 118, 70, 64, 94, 32, 45},
		13: {74, 152, 176, 100, 86, 100, 106, 104, 130, 218, 206, 140, 78},
		15: {8, 183, 61, 91, 202, 37, 51, 58, 58, 237, 140, 124, 5, 99, 105},
		16: {120, 104, 107, 109, 102, 161, 76, 3, 91, 191, 147, 169, 182, 194, 225, 120},
		17: {43, 139, 206, 78, 43, 239, 123, 206, 214, 147, 24, 99, 150, 39, 243, 163, 136},
		18: {215, 234, 158, 94, 184, 97, 118, 170, 79, 187, 152, 148, 252, 179, 5, 98, 96, 153},
		20: {17, 60, 79, 50, 61, 163, 26, 187, 202, 180, 221, 225, 83, 239, 156, 164, 212, 212, 188, 190},
		22: {210, 171, 247, 242, 93, 230, 14, 109, 221, 53, 200, 74, 8, 172, 98, 80, 219, 134, 160, 105, 165, 231},
		24: {229, 121, 135, 48, 211, 117, 251, 126, 159, 180, 169, 152, 192, 226, 228, 218, 111, 0, 117, 232, 87, 96, 227, 21},
		26: {173, 125, 158, 2, 103, 182, 118, 17, 145, 201, 111, 28, 165, 53, 161, 21, 245, 142, 13, 102, 48, 227, 153, 145, 218, 70},
		28: {168, 223, 200, 104, 224, 234, 108, 180, 110, 190, 195, 147, 205, 27, 232, 201, 21, 43, 245, 87, 42, 195, 212, 119, 242, 37, 9, 123},
		30: {41, 173, 145, 152, 216, 31, 179, 182, 50, 48, 110, 86, 239, 96, 222, 125, 42, 173, 226, 193, 224, 130, 156, 37, 251, 216, 238, 40, 192, 180},
	}

	GF = [256]byte{1, 2, 4, 8, 16, 32, 64, 128, 29, 58, 116, 232, 205, 135, 19, 38, 76, 152, 45, 90, 180, 117, 234, 201, 143, 3, 6, 12, 24, 48, 96, 192, 157, 39, 78, 156, 37,
		74, 148, 53, 106, 212, 181, 119, 238, 193, 159, 35, 70, 140, 5, 10, 20, 40, 80, 160, 93, 186, 105, 210, 185, 111, 222, 161,
		95, 190, 97, 194, 153, 47, 94, 188, 101, 202, 137, 15, 30, 60, 120, 240, 253, 231, 211, 187, 107, 214, 177, 127, 254,
		225, 223, 163, 91, 182, 113, 226, 217, 175, 67, 134, 17, 34, 68, 136, 13, 26, 52, 104, 208, 189, 103, 206, 129, 31, 62,
		124, 248, 237, 199, 147, 59, 118, 236, 197, 151, 51, 102, 204, 133, 23, 46, 92, 184, 109, 218, 169, 79, 158, 33, 66, 132,
		21, 42, 84, 168, 77, 154, 41, 82, 164, 85, 170, 73, 146, 57, 114, 228, 213, 183, 115, 230, 209, 191, 99, 198, 145, 63, 126, 252, 229,
		215, 179, 123, 246, 241, 255, 227, 219, 171, 75, 150, 49, 98, 196, 149, 55, 110, 220, 165, 87, 174, 65, 130, 25, 50, 100, 200, 141, 7, 14,
		28, 56, 112, 224, 221, 167, 83, 166, 81, 162, 89, 178, 121, 242, 249, 239, 195, 155, 43, 86, 172, 69, 138, 9, 18, 36, 72, 144, 61, 122, 244, 245,
		247, 243, 251, 235, 203, 139, 11, 22, 44, 88, 176, 125, 250, 233, 207, 131, 27, 54, 108, 216, 173, 71, 142, 1,
	}

	invGF = [256]byte{0, 0, 1, 25, 2, 50, 26, 198, 3, 223, 51, 238, 27, 104, 199, 75,
		4, 100, 224, 14, 52, 141, 239, 129, 28, 193, 105, 248, 200, 8, 76, 113,
		5, 138, 101, 47, 225, 36, 15, 33, 53, 147, 142, 218, 240, 18, 130, 69,
		29, 181, 194, 125, 106, 39, 249, 185, 201, 154, 9, 120, 77, 228, 114, 166,
		6, 191, 139, 98, 102, 221, 48, 253, 226, 152, 37, 179, 16, 145, 34, 136,
		54, 208, 148, 206, 143, 150, 219, 189, 241, 210, 19, 92, 131, 56, 70, 64,
		30, 66, 182, 163, 195, 72, 126, 110, 107, 58, 40, 84, 250, 133, 186, 61,
		202, 94, 155, 159, 10, 21, 121, 43, 78, 212, 229, 172, 115, 243, 167, 87,
		7, 112, 192, 247, 140, 128, 99, 13, 103, 74, 222, 237, 49, 197, 254, 24,
		227, 165, 153, 119, 38, 184, 180, 124, 17, 68, 146, 217, 35, 32, 137, 46,
		55, 63, 209, 91, 149, 188, 207, 205, 144, 135, 151, 178, 220, 252, 190, 97,
		242, 86, 211, 171, 20, 42, 93, 158, 132, 60, 57, 83, 71, 109, 65, 162,
		31, 45, 67, 216, 183, 123, 164, 118, 196, 23, 73, 236, 127, 12, 111, 246,
		108, 161, 59, 82, 41, 157, 85, 170, 251, 96, 134, 177, 187, 204, 62, 90,
		203, 89, 95, 176, 156, 169, 160, 81, 11, 245, 22, 235, 122, 117, 44, 215,
		79, 174, 213, 233, 230, 231, 173, 232, 116, 214, 244, 234, 168, 80, 88, 175,
	}
)
